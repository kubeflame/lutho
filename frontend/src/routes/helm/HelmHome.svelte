<script lang="ts">
  import { type TabItem, type TabQueryParam } from "../../lib/types";
  import HeaderElement from "../../lib/HeaderElement.svelte";
  import Tabs from "../../lib/Tabs.svelte";
  import { tabs } from "../../lib/util";
  import ListReleases from "./ListReleases.svelte";
  import ListRepos from "./ListRepos.svelte";
  import NewRelease from "./NewRelease.svelte";

  const tabItems: TabItem[] = [tabs.list, tabs.repo, tabs.new];

  let tabQueryParam: TabQueryParam;

  $: showNamespaceSelect = tabQueryParam === "list";
</script>

<HeaderElement {showNamespaceSelect}>
  <Tabs
    slot="tabs"
    bind:tabQueryParam
    tabQueryParamDefault={"list"}
    {tabItems}
  />
</HeaderElement>

{#if tabQueryParam === "list"}
  <ListReleases />
{:else if tabQueryParam === "repository"}
  <ListRepos bind:tabQueryParam />
{:else if tabQueryParam === "new"}
  <NewRelease bind:tabQueryParam />
{/if}
