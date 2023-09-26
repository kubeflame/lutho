<script lang="ts">
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import { type TabQueryParam } from "../lib/types";
  import { tabs } from "../lib/util";
  import Dashboard from "../lib/dashboard/Dashboard.svelte";
  import Search from "../lib/Search.svelte";

  const tabItems = [tabs.dashboard, tabs.search];

  let tabQueryParam: TabQueryParam;
</script>

<HeaderElement showNamespaceSelect>
  <Tabs
    slot="tabs"
    bind:tabQueryParam
    tabQueryParamDefault={"dashboard"}
    {tabItems}
  />
</HeaderElement>

{#if tabQueryParam === "dashboard"}
  <Dashboard />
{:else if tabQueryParam === "search"}
  <Search bind:tabQueryParam />
{/if}
