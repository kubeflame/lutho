<script lang="ts">
  import { fade, scale } from "svelte/transition";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import { TabIndex, type KubeGVRK } from "../lib/types";
  import { KubeResourceType, icons, tabs } from "../lib/util";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import Dashboard from "../lib/Dashboard.svelte";
  import Search from "../lib/Search.svelte";
  import RouterPage from "../lib/RouterPage.svelte";

  const tabItems = [tabs.dashboard, tabs.search];

  let isFetching: boolean = false;
  let activeTab: number = TabIndex.DASHBOARD;
  let errorMessage: string;
  let showFilter: boolean;
  let resourceType: KubeGVRK;
  let namespaceInput: string;
  let nameInput: string;
  let labelInput: string;
  let fieldInput: string;
  let onClickSearch: boolean;

  $: if (!showFilter) (labelInput = ""), (fieldInput = "");
</script>

<HeaderElement>
  <Tabs slot="tabs" bind:activeTab {tabItems} />
</HeaderElement>

<RouterPage>
  <ResourceToolbar slot="resource-toolbar" bind:activeTab>
    <div
      slot="homePage"
      class="flex h-8 place-items-center items-center justify-between gap-x-2 rounded-b-2xl
        rounded-t-none border-0 bg-base-100 drop-shadow-md focus:outline-0"
    >
      <button
        on:click={() => (showFilter = !showFilter)}
        class="btn tooltip-primary tooltip btn-sm tooltip-right place-items-center items-center
          rounded-bl-2xl rounded-tl-none bg-base-100 pl-5 font-medium hover:bg-base-200"
        data-tip="more filters"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width={1}
          stroke="currentColor"
          class="h-5 w-5"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d={icons.filter}
          />
        </svg>
      </button>

      <div class="join min-w-fit items-center text-sm" in:scale>
        <label
          for="resource-select"
          class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
        >
          Resource Type
        </label>
        <select
          bind:value={resourceType}
          id="resource-select"
          class="join-item select select-bordered select-xs rounded-lg bg-base-100
            pl-4 text-sm font-normal focus:border-primary/60 focus:outline-0"
        >
          <option selected disabled value="">Choose a resource</option>
          {#each KubeResourceType as resource}
            <option value={resource}>{resource.kind}</option>
          {/each}
        </select>
      </div>

      <div class="join w-full items-center text-sm" in:scale>
        <label
          for="name-input"
          class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
        >
          Name
        </label>
        <input
          id="name-input"
          bind:value={nameInput}
          type="text"
          class="input join-item input-bordered input-xs flex grow
            bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>

      {#if resourceType?.isNamespaced}
        <div class="join w-full items-center text-sm" in:scale>
          <label
            for="namespace-input"
            class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
          >
            Namespace
          </label>
          <input
            id="namespace-input"
            bind:value={namespaceInput}
            type="text"
            class="input join-item input-bordered input-xs flex grow
            bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
          />
        </div>
      {/if}

      {#if showFilter}
        <div class="join w-full items-center text-sm" in:scale>
          <label
            for="field-input"
            class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
          >
            Field Selector
          </label>
          <input
            id="field-input"
            bind:value={fieldInput}
            type="text"
            class="input join-item input-bordered input-xs flex grow
              bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
          />
        </div>

        <div class="join h-6 w-full items-center text-sm" in:scale>
          <label
            for="label-input"
            class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
          >
            Label Selector
          </label>
          <input
            id="label-input"
            bind:value={labelInput}
            type="text"
            class="input join-item input-bordered input-xs flex grow
              bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
          />
        </div>
      {/if}

      <button
        class="btn btn-sm place-items-center items-center rounded-br-2xl
          rounded-tr-none bg-base-100 pr-5 hover:bg-base-200"
        on:click={() => {
          onClickSearch = true;
        }}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width={1}
          stroke="currentColor"
          class="h-5 w-5"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d={icons.arrowRight}
          />
        </svg>
      </button>
    </div>
  </ResourceToolbar>
  {#if isFetching}
    <LoadingNewton />
  {/if}
  {#if errorMessage}
    <ErrorPage {errorMessage} />
  {:else if activeTab === TabIndex.DASHBOARD}
    <Dashboard />
  {:else if activeTab === TabIndex.SEARCH}
    <Search
      bind:isFetching
      bind:errorMessage
      bind:fieldInput
      bind:labelInput
      bind:nameInput
      bind:namespaceInput
      bind:onClickSearch
      {resourceType}
    />
  {/if}
</RouterPage>
