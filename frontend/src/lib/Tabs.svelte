<script lang="ts">
  import { fly } from "svelte/transition";
  import { expoIn } from "svelte/easing";
  import { type TabItem, type TabQueryParam } from "./types";
  import { location, querystring } from "svelte-spa-router";
  import { onMount } from "svelte";
  import SvgIcon from "./SvgIcon.svelte";
  import { writable } from "svelte/store";

  export let tabItems: TabItem[];
  export let activeContainer = writable<string>();
  export let tabQueryParamDefault: TabQueryParam;
  export let tabQueryParam: TabQueryParam;
  export let containers: string[] = [];

  onMount(() => {
    if (!tabQueryParam && !window.location.pathname.startsWith("/srv")) {
      window.location.hash = `${$location}?tab=${tabQueryParamDefault}`;
    }
  });

  $: tabQueryParam = new URLSearchParams($querystring)
    .get("tab")
    ?.toLowerCase() as TabQueryParam;
</script>

<div
  role="tablist"
  class="tabs tabs-lifted mt-3"
  in:fly={{ y: 10, duration: 250, easing: expoIn }}
>
  {#each tabItems as tab}
    <button
      on:click={() => {
        window.location.hash = `${$location}?tab=${tab.name}`;
      }}
      role="tab"
      class="tab space-x-1 {tabQueryParam === tab.name.toLowerCase()
        ? 'tab-active'
        : ''}"
    >
      <SvgIcon strokeWidth={tab.strokeWidth} type={tab.icon} />
      <div>{tab.name}</div>
    </button>
  {/each}
</div>
{#if tabQueryParam === "shell" && containers.length >= 1}
  <div
    class="tooltip tooltip-right tooltip-primary ml-2 mt-1"
    data-tip="available containers"
    in:fly={{ y: 10, duration: 250, easing: expoIn }}
  >
    <select
      bind:value={$activeContainer}
      class="select select-primary select-xs bg-base-100 rounded-full pl-4
        text-center text-sm font-normal drop-shadow-sm focus:outline-1"
    >
      {#each containers as container}
        <option value={container}>{container}</option>
      {/each}
    </select>
  </div>
{:else if tabQueryParam === "logs" && containers.length >= 1}
  <div
    class="tooltip tooltip-right tooltip-primary ml-2 mt-1"
    data-tip="available containers"
    in:fly={{ y: 10, duration: 250, easing: expoIn }}
  >
    <select
      bind:value={$activeContainer}
      class="select select-primary select-xs bg-base-100 rounded-full pl-4
        text-center text-sm font-normal drop-shadow-sm focus:outline-1"
    >
      {#each containers as container}
        <option value={container}>{container}</option>
      {/each}
    </select>
  </div>
{/if}
