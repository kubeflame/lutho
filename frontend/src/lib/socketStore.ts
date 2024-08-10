import { get, writable, type Writable } from "svelte/store";
import {
  type KubeDataStreamMessage,
  type KubeOp,
  type ServerResponse,
  type SockState,
} from "./types";
import { apiURL, httpStatus, WSCloseCode } from "./util";
import SockJS from "sockjs-client/dist/sockjs";
import { onDestroy, onMount } from "svelte";
import { statusCode } from "./stores";

async function getSessionID(
  isLoading: Writable<boolean>,
  sockError: Writable<string>,
  sessionId: Writable<string>,
) {
  isLoading.set(true);

  await fetch(apiURL.dataWS)
    .then((resp) => {
      return resp.json();
    })
    .then((result: ServerResponse) => {
      statusCode.set(result.statusCode);
      sockError.set(result.error);
      sessionId.set(result.sessionId);
      isLoading.set(false);
    })
    .catch((error) => {
      sockError.set(`Cannot get sessionId: ${error}`);
      isLoading.set(false);
    });
}

function sockConnection(
  isLoading: Writable<boolean>,
  sockState: Writable<SockState>,
  sockError: Writable<string>,
  sessionId: Writable<string>,
  dataList: Writable<any>,
  dataGet: Writable<any>,
  dataUpdate: Writable<any>,
  dataDelete: Writable<any>,
): WebSocket {
  const sock: WebSocket = new SockJS(`http://${location.host}${apiURL.data}`);

  getSessionID(isLoading, sockError, sessionId);

  sock.onopen = function () {
    sockState.set({ state: sock.readyState, bound: false });
    sessionId.subscribe((id) => {
      if (id) {
        sock.send(
          JSON.stringify({
            op: { type: "bind" },
            sessionId: id,
          }),
        );
      }
    });
  };

  sock.onmessage = function (me: MessageEvent) {
    const resp: KubeDataStreamMessage = JSON.parse(me.data);
    if (resp.error) {
      statusCode.set(resp.error === httpStatus[401] ? 401 : 502);
      sockError.set(resp.error as string);
      isLoading.set(false);
    } else if (resp.op.type === "bind") {
      sockState.set({ state: sock.readyState, bound: true });
    } else if (resp.op.type === "list") {
      dataList.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "listAll") {
      dataList.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmShowValues") {
      dataGet.set({
        op: resp.op,
        data: resp.data,
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmList") {
      dataList.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "get") {
      dataGet.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmGet") {
      dataGet.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmInstall") {
      dataGet.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmUpgrade") {
      dataGet.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmPull") {
      dataGet.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "helmGetTags") {
      dataGet.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "update") {
      dataUpdate.set({
        op: resp.op,
        data: JSON.parse(resp.data),
      });
      isLoading.set(false);
    } else if (resp.op.type === "delete") {
      dataDelete.set(resp.error);
      isLoading.set(false);
    } else if (resp.op.type === "helmUninstall") {
      dataDelete.set(resp.error);
      isLoading.set(false);
    }
  };

  sock.onerror = function (e: Event) {
    sockError.set(e.toString());
    sock.close(WSCloseCode.error);
    isLoading.set(false);
  };

  sock.onclose = (ce: CloseEvent) => {
    if (!ce.wasClean) {
      sockError.set(ce.reason);
    }
    sockState.subscribe((s) => {
      s.state === WebSocket.OPEN &&
        sock.send(
          JSON.stringify({
            op: { type: "close" },
          }),
        );
    });
    sockState.set({ state: WebSocket.CLOSED, bound: false });
    isLoading.set(false);
  };

  return sock;
}

export default function createSocketStore() {
  const isLoading = writable<boolean>(false);
  const sockState = writable<SockState>({ state: 0, bound: false });
  const sockError = writable<string>();
  const sessionId = writable<string>("");
  const dataSend = writable<KubeOp[]>();
  const dataList = writable<any>();
  const dataGet = writable<any>();
  const dataUpdate = writable<any>();
  const dataDelete = writable<any>();

  let sock: WebSocket;

  onMount(() => {
    sock = sockConnection(
      isLoading,
      sockState,
      sockError,
      sessionId,
      dataList,
      dataGet,
      dataUpdate,
      dataDelete,
    );
  });

  sockState.subscribe((s) => {
    dataSend.subscribe((data) => {
      if (s.state === WebSocket.CLOSED && data && s.refresh && !s.bound) {
        sock = sockConnection(
          isLoading,
          sockState,
          sockError,
          sessionId,
          dataList,
          dataGet,
          dataUpdate,
          dataDelete,
        );
        s.refresh = false;
      } else if (s.state === WebSocket.OPEN && s.bound && data) {
        isLoading.set(true);
        data.forEach((item) => {
          sock.send(
            JSON.stringify({
              op: item,
              sessionId: get(sessionId),
            }),
          );
        });
      }
    });
  });

  onDestroy(() => {
    isLoading.set(false);

    get(sockState).state === WebSocket.OPEN && sock.close(WSCloseCode.info);
  });

  return {
    isLoading,
    dataSend,
    dataList,
    dataGet,
    dataUpdate,
    dataDelete,
    sockError,
    sockState,
  };
}
