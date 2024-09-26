<script lang="ts">
  import SelfSubjectRulesReview from "./SelfSubjectRulesReview.svelte";
  import socketStore from "./socketStore";
  import { authState, kubeHost, namespace } from "./stores";
  import type { Alert, AuthResponse, KubeOp } from "./types";
  import {
    apiURL,
    randomUUID,
    routeString,
    SelfSubjectRulesReviewV1GVRK,
  } from "./util";
  import BottomAlert from "./BottomAlert.svelte";
  import type { V1ResourceRule } from "@kubernetes/client-node";
  import CaptionTag from "./CaptionTag.svelte";
  import { push } from "svelte-spa-router";

  export let dialog: HTMLDialogElement;
  export const onOpen = () => {
    if ($authState && $kubeHost) $dataSend = [sendRulesReview];
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
    if (!$authState && !$kubeHost) return;

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

  async function logout() {
    await fetch(`${location.pathname}${apiURL.auth}`, {
      method: "POST",
      body: JSON.stringify({ type: "authUnset" }),
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((resp) => {
        return resp.json();
      })
      .then((result: AuthResponse) => {
        $authState = result.state;
        $kubeHost = result.kubeHost;
        alert = { message: result.error, type: "error" };
        if (!alert.message) {
          push(routeString.auth);
        }
      })
      .catch((error) => {
        alert = { message: error, type: "error" };
      });
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
      {$kubeHost || "not connected"}
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
          disabled={!$authState && !$kubeHost}
        />
      </div>
      <button
        class="btn btn-primary join-item btn-xs rounded-lg"
        on:click={get}
        disabled={!$authState && !$kubeHost}
      >
        Review
      </button>
    </div>

    <div class="mt-4">
      <SelfSubjectRulesReview rules={rulesData} />
    </div>

    <div class="flex space-x-2">
      <div class="spacer flex grow items-center" />

      {#if $kubeHost && $authState}
        <div class="modal-action flex gap-x-2">
          <button
            class="btn btn-xs outline outline-1 outline-base-100 drop-shadow-md
              hover:bg-warning focus:outline-1 focus:outline-warning"
            on:click={logout}
          >
            Logout
          </button>
        </div>
      {/if}

      <div class="modal-action flex gap-x-2">
        <button
          class="btn btn-xs outline outline-1 outline-base-100 drop-shadow-md
            hover:bg-error focus:outline-1 focus:outline-error"
          on:click={() => dialog.close()}
        >
          Close
        </button>
      </div>
    </div>

    <BottomAlert bind:alert />
  </div>
</dialog>
