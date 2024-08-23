<script lang="ts">
  import { discardYamlChanges } from "./stores";
  import { scale } from "svelte/transition";
  import type { TabQueryParam } from "./types";
  import { writable } from "svelte/store";

  export let tabQueryParam: TabQueryParam = "" as TabQueryParam;
  export let followLogs: boolean = false;
  export let onChangeLogsBtn = writable<Event>();
  export let showShellReconnect = writable<boolean>(false);
  export let reconnectShell = writable<boolean>(false);
  export let onClickSubmit: Function = () => {};
  export let codeMirrorChanged: boolean = false;
</script>

<div
  class="resource-toolbar no-scrollbar bg-base-100 top-10 z-50 flex
    h-8 w-full rounded-b-2xl drop-shadow-md"
>
  <slot name="breadcrumbs" />

  <slot name="custom" />

  <!-- <slot name="helmDefaultValues" /> -->

  {#if tabQueryParam === "yaml" && codeMirrorChanged}
    <div class="spacer flex grow items-center" />

    <div
      class="discard-submit mr-4 flex items-center space-x-2"
      in:scale={{ duration: 200 }}
    >
      <button
        class="btn btn-xs bg-base-100 outline-error hover:btn-error
          text-sm outline outline-1 drop-shadow-md"
        on:click={() => ($discardYamlChanges = true)}
      >
        Discard
      </button>
      <button
        class="btn btn-xs bg-base-100 outline-warning hover:btn-warning
          text-sm outline outline-1 drop-shadow-md"
        on:click={onClickSubmit()}
      >
        Submit
      </button>
    </div>
  {:else if tabQueryParam === "shell" && $showShellReconnect}
    <div class="spacer flex grow items-center" />

    <div
      class="shell-connect flex items-center space-x-2"
      in:scale={{ duration: 200 }}
    >
      <button
        class="btn btn-xs bg-base-100 outline-primary hover:btn-primary
          mr-4 text-sm outline outline-1 drop-shadow-md"
        on:click={() => {
          reconnectShell.set(true);
        }}
      >
        Reconnect
      </button>
    </div>
  {:else if tabQueryParam === "logs"}
    <div class="spacer flex grow items-center" />

    <label class="swap swap-rotate" in:scale={{ duration: 200 }}>
      <input
        type="checkbox"
        bind:checked={followLogs}
        on:change={(ev) => ($onChangeLogsBtn = ev)}
      />
      <div
        class="follow-stop-logs btn btn-xs bg-base-100 mr-4 text-sm outline outline-1 drop-shadow-md
          {followLogs
          ? 'outline-warning hover:btn-warning'
          : 'outline-primary hover:btn-primary'}"
      >
        {followLogs ? "Stop" : "Follow"}
      </div>
    </label>
  {/if}
</div>
