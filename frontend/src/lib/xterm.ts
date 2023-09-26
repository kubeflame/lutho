import SockJS from "sockjs-client/dist/sockjs";
import { Terminal, type ITheme } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { apiURL, colorScheme } from "./util";
import { onDestroy, onMount } from "svelte";
import { writable } from "svelte/store";
import {
  KubeDataOpType,
  WSCloseCode,
  type TerminalMessage,
  type Alert,
  type XtermData,
} from "./types";

const lightTheme: ITheme = {
  background: colorScheme.light["base-100"],
  foreground: colorScheme.light["base-content"],
  cursor: colorScheme.light.primary,
  cursorAccent: colorScheme.light["base-300"],
  selectionBackground: colorScheme.light.primary,
  selectionForeground: colorScheme.light["base-content"],
};

const darkTheme: ITheme = {
  background: colorScheme.dark["base-100"],
  foreground: colorScheme.dark["base-content"],
  cursor: colorScheme.dark.primary,
  cursorAccent: colorScheme.dark["base-300"],
  selectionBackground: colorScheme.dark.primary,
  selectionForeground: colorScheme.dark["base-content"],
};

export default function createXterm(data: XtermData) {
  const element = writable<HTMLElement>();
  const isLoading = writable<boolean>(false);
  const shellReconnect = writable<boolean>();
  const alert = writable<Alert>({ message: null, type: null });
  const sockError = writable<string>();
  const sockState = writable<{ state: number; bound: boolean }>({
    state: 0,
    bound: false,
  });
  const sessionId = writable<string>("");
  const container = writable<string>("");

  const term = new Terminal({
    cursorBlink: true,
    fontFamily: "Fira Code, Iosevka, monospace",
    fontSize: 15,
    cols: 181,
    rows: 33,
  });

  const fitAddon = new FitAddon();
  const xterm_resize_ob = new ResizeObserver(() => {
    fitAddon && fitAddon.fit();
  });

  const sock: WebSocket = new SockJS(`http://${location.host}${apiURL.shell}`);

  onMount(async () => {
    let activeContainer: string = "";

    isLoading.update((load) => (load = true));
    container.subscribe((c) => (activeContainer = c));

    await fetch(
      `${apiURL.shellExec}?namespace=${data.namespace}&name=${data.name}&container=${activeContainer}`,
    )
      .then((resp) => {
        return resp.json();
      })
      .then((result: TerminalMessage) => {
        sockError.update((err) => (err = result.error as string));
        sessionId.update((id) => (id = result.sessionId));
      })
      .catch((error) => {
        sockError.update((err) => (err = `Cannot get sessionId: ${error}`));
      });
  });

  term.options.theme =
    localStorage.getItem("theme") === "light" ? lightTheme : darkTheme;
  term.loadAddon(fitAddon);

  element.subscribe((el) => {
    if (el) {
      term.open(el);
      fitAddon.fit();
      term.focus();

      xterm_resize_ob.observe(el);
    }
  });

  sock.onopen = function () {
    sockState.update((s) => (s = { state: sock.readyState, bound: false }));
    sessionId.subscribe(
      (id) =>
        id &&
        sock.send(
          JSON.stringify({
            op: KubeDataOpType.bind,
            sessionId: id,
          }),
        ),
    );
  };

  sock.onmessage = function (e: any) {
    const resp: TerminalMessage = JSON.parse(e.data);
    if (resp.error) {
      sockError.update((err) => (err = resp.error as string));
    } else if (resp.op === KubeDataOpType.bind) {
      sockState.update((s) => (s = { state: this.readyState, bound: true }));
      isLoading.update((load) => (load = false));
    } else if (resp.op === KubeDataOpType.stdout) {
      term.write(resp.data);
    }
  };

  sock.onerror = function (e: Event) {
    sockError.update((err) => (err = e.toString()));
    sock.close(WSCloseCode.error);
  };

  sock.onclose = (ce: CloseEvent) => {
    isLoading.update((load) => (load = false));
    if (ce.code === WSCloseCode.info) {
      alert.update((a) => (a = { message: ce.reason, type: "info" }));
      shellReconnect.update((sr) => (sr = true));
    } else if (ce.code === WSCloseCode.error) {
      alert.update((a) => (a = { message: ce.reason, type: "error" }));
      shellReconnect.update((sr) => (sr = true));
    } else if (ce.code === WSCloseCode.warning) {
      alert.update((a) => (a = { message: ce.reason, type: "warning" }));
      shellReconnect.update((sr) => (sr = true));
    }

    if (!ce.wasClean) {
      sockError.update((err) => (err = ce.reason));
    }
  };

  sockError.subscribe(
    (err) =>
      err &&
      alert.update((a) => {
        a = { message: err, type: "error" };
        isLoading.update((load) => (load = false));
        return a;
      }),
  );

  term.onResize(({ cols, rows }) => {
    sockState.subscribe((s) => {
      if (s.state === WebSocket.OPEN && s.bound)
        sock.send(
          JSON.stringify({
            op: KubeDataOpType.resize,
            cols: cols,
            rows: rows,
          }),
        );
    });

    fitAddon.fit();
  });

  term.onData((e) => {
    sockState.subscribe((s) => {
      if (s.state === WebSocket.OPEN && s.bound)
        sock.send(
          JSON.stringify({
            op: KubeDataOpType.stdin,
            data: e,
            cols: term.cols,
            rows: term.rows,
          }),
        );
    });
  });

  const htmlEl = document.querySelector("html");
  const callback = (mutationRecords: MutationRecord[]) => {
    for (const record of mutationRecords) {
      if (record.type === "attributes") {
        if (record.attributeName === "data-theme") {
          localStorage.getItem("theme") === "light"
            ? (term.options.theme = lightTheme)
            : (term.options.theme = darkTheme);
        }
      }
    }
  };
  const observer = new MutationObserver(callback);
  observer.observe(htmlEl as Node, { attributes: true });

  onDestroy(() => {
    sockState.subscribe((s) => {
      s.state === WebSocket.OPEN && sock.close(WSCloseCode.info);
    });
    term.dispose();
  });

  return {
    container,
    element,
    isLoading,
    alert,
    shellReconnect,
    sockError,
    term,
  };
}
