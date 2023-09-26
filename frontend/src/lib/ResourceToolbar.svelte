<script lang="ts">
  import { link } from "svelte-spa-router";
  import { discardYamlChanges } from "./stores";
  import { icons } from "./util";
  import { fly, scale } from "svelte/transition";
  import { type ResourceToolbarContent, TabIndex } from "./types";

  export let toolbarContent: ResourceToolbarContent[] = [];
  export let activeTab: number = TabIndex.DETAILS;
  export let followLogs: boolean = false;
  export let stopLogs: boolean = false;
  export let showShellReconnect: boolean = false;
  export let refreshShellComponent: boolean = false;
  export let onClickSubmit: Function = () => {};
  export let codeMirrorChanged: boolean = false;
</script>

{#if toolbarContent.length >= 1}
  <div
    class="resource-toolbar sticky top-10 -z-10 flex h-8 transform items-center gap-2
  rounded-b-2xl bg-base-100 px-5 drop-shadow-md duration-300 ease-in-out"
    in:fly={{ duration: 200, y: -1 }}
  >
    <div class="breadcrumbs text-sm">
      <ul>
        <div class="flex pl-0.5 pr-1.5">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width={1}
            stroke="currentColor"
            class="h-4 w-4"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d={icons.at} />
          </svg>
        </div>
        {#each toolbarContent as item}
          <li in:scale={{ duration: 200 }}>
            {#if item.url}
              <a
                class="drop-shadow-2xl"
                href={item.url}
                use:link
                on:click={() => {
                  toolbarContent = toolbarContent.filter((x) => {
                    return x.index <= item.index;
                  });
                  item.url = "";
                }}
              >
                {item.name}
              </a>
            {:else}
              {item.name}
            {/if}
          </li>
        {/each}
      </ul>
    </div>

    <div class="spacer flex grow items-center" />

    {#if activeTab === TabIndex.YAML && codeMirrorChanged}
      <div class="discard-submit space-x-2" in:scale={{ duration: 200 }}>
        <button
          class="btn btn-xs bg-base-100 text-sm outline outline-1 outline-error drop-shadow-md hover:btn-error"
          on:click={() => ($discardYamlChanges = true)}
        >
          Discard
        </button>
        <button
          class="btn btn-xs bg-base-100 text-sm outline outline-1 outline-warning drop-shadow-md hover:btn-warning"
          on:click={onClickSubmit()}
        >
          Submit
        </button>
      </div>
    {:else if activeTab === TabIndex.SHELL && showShellReconnect}
      <div class="shell-connect space-x-2" in:scale={{ duration: 200 }}>
        <button
          class="btn btn-xs bg-base-100 text-sm outline outline-1 outline-primary drop-shadow-md hover:btn-primary"
          on:click={() => {
            refreshShellComponent = true;
          }}
        >
          Reconnect
        </button>
      </div>
    {:else if activeTab === TabIndex.LOGS}
      <label class="swap swap-rotate" in:scale={{ duration: 200 }}>
        <input
          type="checkbox"
          bind:checked={followLogs}
          on:change={() => {
            stopLogs = !followLogs;
          }}
        />
        <div
          class="follow-stop-logs btn btn-xs bg-base-100 text-sm outline outline-1 drop-shadow-md
            {followLogs
            ? 'outline-warning hover:btn-warning'
            : 'outline-primary hover:btn-primary'}"
        >
          {followLogs ? "Stop" : "Follow"}
        </div>
      </label>
    {/if}
  </div>
{:else if activeTab === TabIndex.SEARCH}
  <slot name="homePage" />
{:else}
  <div class="flex h-8 rounded-b-2xl bg-base-100 drop-shadow-md" />
{/if}
