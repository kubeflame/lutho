<script lang="ts">
  import SelfSubjectRulesReview from "./SelfSubjectRulesReview.svelte";
  import socketStore from "./socketStore";
  import { kubeHost, namespace } from "./stores";
  import type { Alert, KubeOp } from "./types";
  import { randomUUID, SelfSubjectRulesReviewV1GVRK } from "./util";
  import BottomAlert from "./BottomAlert.svelte";
  import type { V1ResourceRule } from "@kubernetes/client-node";
  import CaptionTag from "./CaptionTag.svelte";

  export let dialog: HTMLDialogElement;
  export const onOpen = () => {
    $dataSend = [sendRulesReview];
  };

  const { sockState, sockError, isLoading, dataSend, dataGet } = socketStore();

  let rulesData: Array<V1ResourceRule>;
  let rulesNamespace = "";
  let alert: Alert = {
    message: null,
    type: null,
    className: "rounded-lg mt-4 p-1 pl-2",
  };

  $: !rulesNamespace && (rulesNamespace = $namespace);
  $: sendRulesReview = {
    opID: randomUUID(),
    type: "rulesReview",
    request: {
      kubeGVRK: SelfSubjectRulesReviewV1GVRK,
      namespace: rulesNamespace || $namespace,
    },
  } as KubeOp;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendRulesReview.opID) {
      rulesData = dg.data?.status?.resourceRules;
    }
  });

  namespace.subscribe((n) => (rulesNamespace = n));

  function get() {
    resetData();

    if ($sockState.state === WebSocket.CLOSED) {
      $sockState.refresh = true;
    } else {
      $dataSend = [sendRulesReview];

      sockError.subscribe((err) => {
        if (err) {
          alert.message = err;
          alert.type = "error";
        }
      });
    }
  }

  function resetData() {
    alert.message = null;
    alert.type = null;
    rulesData = [];
    $dataGet = "";
    $sockError = "";
  }
</script>

<dialog bind:this={dialog} class="modal modal-middle" on:close={resetData}>
  <div
    class="modal-box grid max-h-[calc(70vh)] max-w-[calc(50vw)]
      grid-cols-1 rounded-lg border border-base-200 bg-base-100
      p-2 outline outline-1 outline-base-200 drop-shadow-md"
  >
    <CaptionTag tagName={"Cluster connection"} bgClassNames={"mt-2"} />
    <div class="font-mono">
      {$kubeHost}
    </div>

    <CaptionTag tagName={"Namespace access rules"} bgClassNames={"mt-4"} />
    <div class="join flex h-max w-max space-x-2 rounded-md">
      <div class="join join-item">
        <label
          for="namespace"
          class="join-item h-full bg-base-200 p-0.5 pl-2 pr-2 text-sm font-normal"
        >
          Namespace
        </label>
        <input
          id="namespace"
          bind:value={rulesNamespace}
          type="text"
          class="input input-xs join-item input-bordered flex grow
            bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>
      <button
        class="btn btn-primary join-item btn-xs rounded-lg"
        on:click={get}
      >
        Review
      </button>
    </div>

    <div class="mt-4">
      <SelfSubjectRulesReview rules={rulesData} />
    </div>

    <div class="modal-action flex gap-x-2">
      <button
        class="btn btn-xs outline outline-1 outline-base-100 drop-shadow-md
          hover:bg-error focus:outline-1 focus:outline-error"
        on:click={() => dialog.close()}
      >
        Close
      </button>
    </div>

    <BottomAlert bind:alert />
  </div>
</dialog>
