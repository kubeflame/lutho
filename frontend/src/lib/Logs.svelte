<script lang="ts">
  import type { Alert, LogStreamMessage } from "./types";
  import { onDestroy } from "svelte";
  import { apiURL, transitionEffects, WSCloseCode } from "./util";
  import Convert from "ansi-to-html";
  import BottomAlert from "./BottomAlert.svelte";
  import LoadingNewton from "./LoadingNewton.svelte";
  import type { Writable } from "svelte/store";

  export let name: string = "";
  export let namespace: string = "";
  export let followLogs: boolean = false;
  export let activeContainer: Writable<string>;
  export let onChangeLogsBtn: Writable<Event>;

  const convert = new Convert({ newline: true });

  let isFetching: boolean;
  let sessionId = "";
  let sock: WebSocket;
  let el: Element;

  onChangeLogsBtn.subscribe((c) => {
    if (c && !followLogs) {
      sock?.close(WSCloseCode.warning, "Log tailing stopped");
      alert = { message: "Log tailing stopped", type: "warning" };
    }
  });

  $: $activeContainer && getLogsBasedOnContainer($activeContainer);
  $: alert = { message: null, type: null } as Alert;
  $: if (followLogs && $activeContainer) {
    el?.insertAdjacentHTML(
      "beforeend",
      "<div class='divider text-info'>   ** Log tailing started **   </div>",
    );
    newSession($activeContainer);
  }

  function getLogsBasedOnContainer(ac: string) {
    if (followLogs && sock?.readyState === WebSocket.OPEN) {
      followLogs = false;
      sock?.close(WSCloseCode.warning, "Log tailing stopped.");
    }
    newSession(ac);
  }

  async function newSession(container: string) {
    isFetching = true;
    alert = { message: null, type: null };

    if (el && !followLogs) el.innerHTML = "";

    await fetch(
      `${apiURL.logsStream}?namespace=${namespace}&name=${
        name
      }&container=${container}&follow=${followLogs}&tailLines=${
        followLogs ? "1" : "1000"
      }`,
    )
      .then((resp) => {
        return resp.json();
      })
      .then((result: LogStreamMessage) => {
        sessionId = result.sessionId;
      })
      .catch((error) => {
        alert = { message: error, type: "error" };
      });

    sock = new WebSocket(`ws://${location.host}${apiURL.logs}`);

    sock.onopen = function () {
      isFetching = false;
      const startData = { Op: "bind", SessionID: sessionId };
      sock.send(JSON.stringify(startData));
    };

    sock.onmessage = function (e: MessageEvent<any>) {
      let resp: LogStreamMessage = JSON.parse(e.data);
      if (resp.op === "bind") {
        sock.send(
          JSON.stringify({
            op: "stdin",
            sessionId: sessionId,
          }),
        );
      } else if (resp.op === "close") {
        if (resp.statusCode === WSCloseCode.info) {
          alert = { message: resp.data, type: "info" };
        } else if (
          resp.statusCode === WSCloseCode.warning &&
          resp.data === "EOF"
        ) {
          alert = { message: resp.data, type: "warning" };
        } else if (resp.statusCode === WSCloseCode.error) {
          alert = { message: resp.data, type: "error" };
        } else if (resp.statusCode === WSCloseCode.warning) {
          alert = { message: resp.data, type: "warning" };
        }
      } else if (resp.op === "stdout") {
        el?.insertAdjacentHTML("beforeend", convert.toHtml(resp.data));
        if (followLogs) {
          el.scroll({
            top: el.scrollHeight,
            behavior: "auto",
          });
        }
      }
    };

    sock.onclose = function (ce) {};

    sock.onerror = function (e) {
      alert = { message: e.toString(), type: "error" };
      sock?.readyState === WebSocket.OPEN && sock?.close(WSCloseCode.error);
    };
  }

  onDestroy(() => {
    followLogs = false;
    sock?.readyState === WebSocket.OPEN && sock?.close(WSCloseCode.info);
  });
</script>

{#if isFetching}
  <LoadingNewton />
{/if}

<div
  class="logs-element flex h-full w-full rounded-lg font-mono text-sm {transitionEffects}"
>
  <div
    bind:this={el}
    class="logs h-full w-full overflow-y-scroll p-4 {transitionEffects}"
  />
  <BottomAlert bind:alert />
</div>
