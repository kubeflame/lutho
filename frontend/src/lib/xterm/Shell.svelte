<script lang="ts">
  import "@xterm/xterm/css/xterm.css";
  import BottomAlert from "../BottomAlert.svelte";
  import LoadingNewton from "../LoadingNewton.svelte";
  import { apiURL, transitionEffects, WSCloseCode } from "../util";
  import { type Writable } from "svelte/store";
  import type { Alert, TerminalMessage } from "../types";
  import { Terminal } from "@xterm/xterm";
  import { FitAddon } from "@xterm/addon-fit";
  import { darkTheme, lightTheme } from "./themes";
  import { onDestroy } from "svelte";
  import SockJS from "sockjs-client/dist/sockjs";

  export let name: any;
  export let namespace: any;
  export let showShellReconnect: Writable<boolean>;
  export let activeContainer: Writable<string>;
  export let reconnectShell: Writable<boolean>;

  let isLoading = false;
  let errorMessage = "";
  let el: HTMLElement;
  let sock: WebSocket;
  let sessionId = "";
  let term: Terminal;
  let sockState = { state: 0, bound: false };
  let alert = { message: null, type: null } as Alert;

  $: $activeContainer && el && init(el);
  $: errorMessage && (alert = { message: errorMessage, type: "error" });
  $: $reconnectShell && init(el);

  function init(el: any) {
    createSocket();
    createXterm(el);
  }

  async function createSocket() {
    if (sock) sock.close();
    alert = { message: null, type: null };
    reconnectShell.set(false);
    isLoading = true;

    await fetch(
      `${apiURL.shellExec}?namespace=${namespace}&name=${name}&container=${$activeContainer}`,
    )
      .then((resp) => {
        return resp.json();
      })
      .then((result: TerminalMessage) => {
        errorMessage = result.error as string;
        sessionId = result.sessionId;
      })
      .catch((error) => {
        errorMessage = `Cannot get sessionId: ${error}`;
      });

    sock = new SockJS(`http://${location.host}${apiURL.shell}`);

    sock.onopen = function () {
      showShellReconnect.set(false);
      sockState = { state: this.readyState, bound: false };
      sock.send(
        JSON.stringify({
          op: "bind",
          sessionId: sessionId,
        }),
      );
    };

    sock.onmessage = function (e: any) {
      const resp: TerminalMessage = JSON.parse(e.data);
      if (resp.error) {
        isLoading = false;
        errorMessage = resp.error as string;
      } else if (resp.op === "bind") {
        sockState = { state: this.readyState, bound: true };
        isLoading = false;
        sock.send(
          JSON.stringify({
            op: "resize",
            cols: term.cols,
            rows: term.rows,
          }),
        );
      } else if (resp.op === "stdout") {
        isLoading = false;
        term.write(resp.data);
      }
    };

    sock.onerror = function (e: Event) {
      errorMessage = e.toString();
      sock.readyState === WebSocket.OPEN && sock.close(WSCloseCode.error);
    };

    sock.onclose = (ce: CloseEvent) => {
      isLoading = false;
      showShellReconnect.set(true);
      if (ce.code === WSCloseCode.info) {
        alert = { message: ce.reason, type: "info" };
      } else if (ce.code === WSCloseCode.error) {
        alert = { message: ce.reason, type: "error" };
      } else if (ce.code === WSCloseCode.warning) {
        alert = { message: ce.reason, type: "warning" };
      }

      if (!ce.wasClean) {
        errorMessage = ce.reason;
      }
    };
  }

  function createXterm(element: HTMLElement) {
    if (term) term.dispose();
    term = new Terminal({
      cursorBlink: true,
      fontFamily: "Fira Code, Iosevka, monospace",
      fontSize: 15,
      scrollback: 1000,
    });

    const fitAddon = new FitAddon();
    const xterm_resize_ob = new ResizeObserver(() => {
      fitAddon && fitAddon.fit();
    });
    term.loadAddon(fitAddon);

    term.options.theme =
      localStorage.getItem("theme") === "light" ? lightTheme : darkTheme;

    term.open(element);
    fitAddon.fit();
    term.focus();

    xterm_resize_ob.observe(element);

    term.onResize(({ cols, rows }) => {
      if (sock.readyState === WebSocket.OPEN && sockState.bound) {
        sock.send(
          JSON.stringify({
            op: "resize",
            cols: cols,
            rows: rows,
          }),
        );
      }

      fitAddon.fit();
    });

    term.onData((e) => {
      if (sock.readyState === WebSocket.OPEN && sockState.bound) {
        sock.send(
          JSON.stringify({
            op: "stdin",
            data: e,
            cols: term.cols,
            rows: term.rows,
          }),
        );
      }
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
  }

  onDestroy(() => {
    term && term.dispose();
    sock?.readyState === WebSocket.OPEN && sock.close(WSCloseCode.info);
  });
</script>

{#if isLoading}
  <LoadingNewton />
{/if}

<div bind:this={el} class="terminal h-full w-full {transitionEffects}">
  <BottomAlert bind:alert />
</div>
