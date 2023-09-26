<script lang="ts">
  import { NamespaceV1GVRK } from "./util";
  import type { V1NamespaceList } from "@kubernetes/client-node/dist/gen/api";
  import { namespace } from "./stores";
  import { fly } from "svelte/transition";
  import { expoIn } from "svelte/easing";
  import { KubeDataOpType } from "./types";
  import socketStore from "./socketStore";

  const { sockError, isLoading, dataSend, dataList } = socketStore();

  let nsListData: V1NamespaceList;

  $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...NamespaceV1GVRK,
      },
    },
  ];

  $: nsListData = $dataList;
</script>

{#if $sockError}
  <div
    class="namespace-tooltip tooltip tooltip-error tooltip-bottom mt-1"
    data-tip={$sockError}
  >
    <select
      id="namespace-select"
      class="namespace-select select select-error select-xs rounded-full bg-base-100 text-center text-sm font-normal drop-shadow-sm focus:outline-1"
      in:fly={{ y: -10, duration: 300, easing: expoIn }}
    >
      <option disabled selected>not available</option>
    </select>
  </div>
{:else if nsListData}
  <select
    id="namespace-select"
    bind:value={$namespace}
    class="namespace-select select select-primary select-xs mt-1 rounded-full bg-base-100 text-center text-sm font-normal drop-shadow-sm focus:outline-1"
    in:fly={{ y: -10, duration: 300, easing: expoIn }}
    on:change={() => {
      localStorage.setItem("namespace", $namespace);
    }}
  >
    {#each nsListData.items as ns}
      <option value={ns.metadata?.name}>{ns.metadata?.name}</option>
    {/each}
    <option value="">all namespaces</option>
  </select>
{/if}
