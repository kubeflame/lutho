<script lang="ts">
  import SockJS from "sockjs-client/dist/sockjs";
  import { WSCloseCode, type Alert, type LogStreamMessage } from "./types";
  import { onDestroy } from "svelte";
  import { availableContainers } from "./stores";
  import { apiURL } from "./util";
  import Convert from "ansi-to-html";
  import BottomAlert from "./BottomAlert.svelte";
  import LoadingNewton from "./LoadingNewton.svelte";

  export let termParams: any;
  export let followLogs: boolean;
  export let stopLogs: boolean;
  export let refreshLogsComponent: boolean;
  export let activeContainer: string = "";

  const convert = new Convert();

  $: alert = { message: null, type: null } as Alert;

  let isFetching: boolean;
  let sessionId = "";
  let sock: WebSocket;
  let el: Element;

  async function newSession() {
    isFetching = true;
    alert = { message: null, type: null };

    $availableContainers = termParams.containers;
    if (activeContainer === "") {
      activeContainer = $availableContainers[0];
    }

    sock = new SockJS(`http://${location.host}${apiURL.logs}`);

    sessionId = await fetch(
      `${apiURL.logsStream}?namespace=${termParams.namespace}&name=${
        termParams.name
      }&container=${activeContainer}&follow=${followLogs}&tailLines=${
        followLogs ? "1" : ""
      }`,
    )
      .then((resp) => {
        return resp.json();
      })
      .then((result: LogStreamMessage) => {
        return result.sessionId;
      });

    sock.onopen = function () {
      isFetching = false;
      const startData = { Op: "bind", SessionID: sessionId };
      sock.send(JSON.stringify(startData));
    };

    sock.onmessage = function (e: MessageEvent<any>) {
      let resp: LogStreamMessage = JSON.parse(e.data);
      if (resp.op === "stdout") {
        el.insertAdjacentHTML("beforeend", convert.toHtml(resp.data) + "</br>");
        if (followLogs) {
          el.scroll({
            top: el.scrollHeight,
            behavior: "auto",
          });
        }
      }
    };

    sock.onclose = function (e) {
      const endData = { Op: "close", SessionID: sessionId };
      sock.send(JSON.stringify(endData));
      refreshLogsComponent = false;
      if (e.code === WSCloseCode.info) {
        alert = { message: e.reason, type: "info" };
      } else if (e.code === WSCloseCode.error && e.reason === "EOF") {
        alert = { message: e.reason, type: "warning" };
      } else if (e.code === WSCloseCode.error) {
        alert = { message: e.reason, type: "error" };
      } else if (e.code === WSCloseCode.warning) {
        alert = { message: e.reason, type: "warning" };
      }
    };

    sock.onerror = function (e) {
      console.error(e);
    };
  }

  newSession();

  $: if (refreshLogsComponent) {
    el.innerHTML = "";
    newSession();
  }

  $: if (followLogs) {
    el.insertAdjacentHTML(
      "beforeend",
      "<div class='divider text-info'>   ** Log tailing started **   </div>",
    );
    newSession();
  }

  $: if (stopLogs) {
    const endData = { Op: "close", SessionID: sessionId };
    sock.send(JSON.stringify(endData));
    sock.close(WSCloseCode.warning, "Log tailing stopped");
    stopLogs = false;
  }

  onDestroy(async () => {
    followLogs = false;
    stopLogs = false;
    sock.close(WSCloseCode.info);
  });
</script>

{#if isFetching}
  <LoadingNewton />
{/if}

<div
  class="logs-element flex h-full w-full transform rounded-lg font-mono text-sm duration-300 ease-in-out"
>
  <div bind:this={el} class="logs h-full w-full overflow-y-scroll p-4" />
  <BottomAlert bind:alert />
</div>
