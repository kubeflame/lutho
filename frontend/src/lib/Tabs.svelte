<script lang="ts">
  import { fly } from "svelte/transition";
  import { expoIn } from "svelte/easing";
  import { availableContainers, showUtils } from "./stores";
  import { TabIndex, type TabItem } from "./types";

  export let tabItems: Array<TabItem>;
  export let activeTab: number;
  export let refreshLogsComponent: boolean = false;
  export let refreshShellComponent: boolean = false;
  export let activeContainer: string = "";
</script>

<div class="tabs mt-2" in:fly={{ y: 10, duration: 250, easing: expoIn }}>
  {#each tabItems as tab}
    <button
      on:click={() => {
        activeTab = tab.index;
      }}
      class="tab tab-lifted space-x-1 {activeTab === tab.index
        ? 'tab-active'
        : ''}"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width={tab.strokeWidth}
        stroke="currentColor"
        class="h-4 w-4"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d={tab.icon} />
      </svg>
      <div>{tab.name}</div>
    </button>
  {/each}
</div>
{#if activeTab === TabIndex.SHELL && $availableContainers.length >= 1}
  <div
    class="tooltip-primary tooltip tooltip-right ml-2 mt-1"
    data-tip="available containers"
    in:fly={{ y: 10, duration: 250, easing: expoIn }}
  >
    <select
      bind:value={activeContainer}
      class="select select-primary select-xs rounded-full bg-base-100 pl-4 text-center text-sm font-normal drop-shadow-sm focus:outline-1"
      on:change={() => {
        refreshShellComponent = true;
      }}
    >
      {#each $availableContainers as container}
        <option value={container}>{container}</option>
      {/each}
    </select>
  </div>
{:else if activeTab === TabIndex.LOGS && $availableContainers.length >= 1}
  <div
    class="tooltip-primary tooltip tooltip-right ml-2 mt-1"
    data-tip="available containers"
    in:fly={{ y: 10, duration: 250, easing: expoIn }}
  >
    <select
      bind:value={activeContainer}
      class="select select-primary select-xs rounded-full bg-base-100 pl-4 text-center text-sm font-normal drop-shadow-sm focus:outline-1"
      on:change={() => {
        refreshLogsComponent = true;
      }}
    >
      {#each $availableContainers as container}
        <option value={container}>{container}</option>
      {/each}
    </select>
  </div>
{:else if activeTab === TabIndex.YAML}
  <div
    class="tooltip-primary tooltip tooltip-right ml-2 mt-1"
    in:fly={{ y: 10, duration: 250, easing: expoIn }}
  >
    <button
      class="btn btn-xs rounded-full font-normal outline outline-1 outline-base-100 drop-shadow-sm hover:bg-primary hover:shadow-sm"
      on:click={() => {
        $showUtils = true;
      }}
    >
      Utilities
    </button>
  </div>
{/if}
