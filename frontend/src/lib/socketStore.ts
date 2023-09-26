import { writable } from "svelte/store";
import {
  KubeDataOpType,
  WSCloseCode,
  type KubeDataStreamMessage,
  type KubeOp,
  type ServerResponse,
} from "./types";
import { apiURL } from "./util";
import SockJS from "sockjs-client/dist/sockjs";
import { onDestroy, onMount } from "svelte";
import { statusCode } from "./stores";

export default function createSocketStore() {
  const isLoading = writable<boolean>(false);
  const dataSend = writable<KubeOp[]>();
  const dataList = writable<any>();
  const dataGet = writable<any>();
  const dataUpdate = writable<any>();
  const dataDelete = writable<any>();
  const sockError = writable<string>();
  const sockState = writable<{
    state: number;
    bound: boolean;
  }>({
    state: 0,
    bound: false,
  });

  const sessionId = writable<string>("");

  const sock: WebSocket = new SockJS(`http://${location.host}${apiURL.data}`);

  onMount(async () => {
    isLoading.update((load) => (load = true));

    await fetch(apiURL.dataWS)
      .then((resp) => {
        return resp.json();
      })
      .then((result: ServerResponse) => {
        statusCode.update((sc) => (sc = result.statusCode));
        sockError.update((err) => (err = result.error));
        sessionId.update((id) => (id = result.sessionId));
        isLoading.update((load) => (load = false));
      })
      .catch((error) => {
        sockError.update((err) => (err = `Cannot get sessionId: ${error}`));
        isLoading.update((load) => (load = false));
      });
  });

  sock.onopen = function () {
    sockState.update((s) => (s = { state: sock.readyState, bound: false }));
    sessionId.subscribe(
      (id) =>
        id &&
        sock.send(
          JSON.stringify({
            op: { type: KubeDataOpType.bind },
            sessionId: id,
          }),
        ),
    );
  };

  sock.onmessage = function (me: MessageEvent) {
    const resp: KubeDataStreamMessage = JSON.parse(me.data);
    if (resp.error) {
      sockError.update((err) => (err = resp.error as string));
    } else if (resp.op.type === KubeDataOpType.bind) {
      sockState.update((s) => (s = { state: this.readyState, bound: true }));
    } else if (resp.op.type === KubeDataOpType.list) {
      dataList.update((val) => (val = JSON.parse(resp.data)));
    } else if (resp.op.type === KubeDataOpType.get) {
      dataGet.update((val) => (val = JSON.parse(resp.data)));
    } else if (resp.op.type === KubeDataOpType.update) {
      dataUpdate.update((val) => (val = JSON.parse(resp.data)));
    } else if (resp.op.type === KubeDataOpType.delete) {
      dataDelete.update((val) => (val = resp.error));
    }
    isLoading.update((load) => (load = false));
  };

  sock.onerror = function (e: Event) {
    sockError.update((err) => (err = e.toString()));
    sock.close(WSCloseCode.error);
  };

  sock.onclose = (ce: CloseEvent) => {
    if (!ce.wasClean) {
      sockError.update((err) => (err = ce.reason));
    }
    sockState.subscribe((s) => {
      s.state === WebSocket.OPEN &&
        sock.send(
          JSON.stringify({
            op: { type: KubeDataOpType.close },
          }),
        );
    });
  };

  dataSend.subscribe((data) => {
    sockState.subscribe((s) => {
      if (s.state === WebSocket.OPEN && s.bound && data) {
        isLoading.update((load) => (load = true));
        data.forEach((item) =>
          sock.send(
            JSON.stringify({
              op: item,
            }),
          ),
        );
      }
    });
  });

  onDestroy(() => {
    sockState.subscribe((s) => {
      s.state === WebSocket.OPEN && sock.close(WSCloseCode.info);
    });
  });

  return {
    isLoading,
    dataSend,
    dataList,
    dataGet,
    dataUpdate,
    dataDelete,
    sockError,
  };
}
