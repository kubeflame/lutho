<script lang="ts">
  import { NamespaceV1GVRK, randomUUID } from "./util";
  import type { V1NamespaceList } from "@kubernetes/client-node/dist/gen/api";
  import { namespace } from "./stores";
  import socketStore from "./socketStore";

  const { sockError, dataSend, dataList } = socketStore();

  let nsListData: V1NamespaceList;

  $: sendNsList = {
    opID: randomUUID(),
    type: "list",
    request: {
      kubeGVRK: NamespaceV1GVRK,
    },
  } as any;

  $: $dataSend = [sendNsList];

  $: nsStyle = $sockError
    ? "!text-base-content/55 border !border-error bg-base-200 h-full flex place-items-center"
    : "border !border-base-300 bg-base-100 h-full flex text-center place-items-center";

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendNsList.opID) {
      nsListData = dl.data;
    }
  });
</script>

{#if $sockError}
  <div
    class="namespace-tooltip tooltip tooltip-bottom tooltip-error join flex !h-7
      items-center rounded-xl !shadow-sm"
    data-tip={$sockError}
  >
    <label
      for="namespace-select"
      class="join-item min-w-fit max-w-fit pl-2.5 pr-2 font-mono font-thin tracking-wide {nsStyle}"
    >
      NS
    </label>
    <select
      id="namespace-select"
      class="join-item namespace-select select select-xs rounded-xl
        font-mono text-sm focus:outline-none {nsStyle}"
      disabled
    >
      <option>not available</option>
    </select>
  </div>
{:else if nsListData}
  <div class="join flex !h-7 items-center rounded-xl !shadow-sm">
    <label
      for="namespace-select"
      class="join-item min-w-fit max-w-fit pl-2.5 pr-2 font-mono font-thin tracking-wide {nsStyle}"
    >
      NS
    </label>
    <select
      id="namespace-select"
      bind:value={$namespace}
      class="join-item namespace-select select select-xs rounded-xl
        font-mono text-sm focus:outline-none {nsStyle}"
      on:change={() => {
        localStorage.setItem("namespace", $namespace);
      }}
    >
      {#each nsListData.items as ns}
        <option value={ns.metadata?.name}>{ns.metadata?.name}</option>
      {/each}
      <option value="">all namespaces</option>
    </select>
  </div>
{/if}
