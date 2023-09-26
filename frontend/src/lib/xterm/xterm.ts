import SockJS from "sockjs-client/dist/sockjs";
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { apiURL, WSCloseCode } from "../util";
import { onDestroy, onMount } from "svelte";
import { writable, type Writable } from "svelte/store";
import { type TerminalMessage, type Alert, type SockState } from "../types";
import { darkTheme, lightTheme } from "./themes";

async function getSessionID(
  name: string,
  namespace: string,
  container: string,
  sockError: Writable<string>,
  sessionId: Writable<string>,
) {
  await fetch(
    `${apiURL.shellExec}?namespace=${namespace}&name=${name}&container=${container}`,
  )
    .then((resp) => {
      return resp.json();
    })
    .then((result: TerminalMessage) => {
      sockError.set(result.error as string);
      sessionId.set(result.sessionId);
    })
    .catch((error) => {
      sockError.set(`Cannot get sessionId: ${error}`);
    });
}

export default function createXterm(
  name: string,
  namespace: string,
  activeContainer: Writable<string>,
  showShellReconnect: Writable<boolean>,
) {
  const element = writable<HTMLElement>();
  const isLoading = writable<boolean>(false);
  const alert = writable<Alert>({ message: null, type: null });
  const sockError = writable<string>();
  const sockState = writable<SockState>({ state: 0, bound: false });
  const sessionId = writable<string>("");

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

  onMount(() => {
    isLoading.set(true);
    activeContainer.subscribe((ac) => {
      if (ac) getSessionID(name, namespace, ac, sockError, sessionId);
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
    showShellReconnect.set(false);
    sockState.set({ state: this.readyState, bound: false });
    sessionId.subscribe((id) => {
      if (id) {
        console.log(id);
        
        sock.send(
          JSON.stringify({
            op: "bind",
            sessionId: id,
          }),
        );
      }
    });
  };

  sock.onmessage = function (e: any) {
    const resp: TerminalMessage = JSON.parse(e.data);
    if (resp.error) {
      isLoading.set(false);
      sockError.set(resp.error as string);
    } else if (resp.op === "bind") {
      sockState.set({ state: this.readyState, bound: true });
      isLoading.set(false);
    } else if (resp.op === "stdout") {
      isLoading.set(false);
      term.write(resp.data);
    }
  };

  sock.onerror = function (e: Event) {
    sockError.set(e.toString());
    sock.readyState === WebSocket.OPEN && sock.close(WSCloseCode.error);
  };

  sock.onclose = (ce: CloseEvent) => {
    isLoading.set(false);
    showShellReconnect.set(true);
    if (ce.code === WSCloseCode.info) {
      alert.set({ message: ce.reason, type: "info" });
    } else if (ce.code === WSCloseCode.error) {
      alert.set({ message: ce.reason, type: "error" });
    } else if (ce.code === WSCloseCode.warning) {
      alert.set({ message: ce.reason, type: "warning" });
    }

    if (!ce.wasClean) {
      sockError.set(ce.reason);
    }
  };

  sockError.subscribe((err) => {
    if (err) {
      alert.set({ message: err, type: "error" });
      isLoading.set(false);
    }
  });

  term.onResize(({ cols, rows }) => {
    if (sock.readyState === WebSocket.OPEN)
      sock.send(
        JSON.stringify({
          op: "resize",
          cols: cols,
          rows: rows,
        }),
      );

    fitAddon.fit();
  });

  term.onData((e) => {
    if (sock.readyState === WebSocket.OPEN)
      sock.send(
        JSON.stringify({
          op: "stdin",
          data: e,
          cols: term.cols,
          rows: term.rows,
        }),
      );
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
    term.dispose();
    sock.readyState === WebSocket.OPEN && sock?.close(WSCloseCode.info);
  });

  return {
    element,
    isLoading,
    alert,
    sockError,
    term,
  };
}
