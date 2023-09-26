<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { routeString, tabs } from "../lib/util";
  import { TabIndex, type Helm, KubeDataOpType } from "../lib/types";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Details from "../lib/Details.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let activeTab: number = TabIndex.DETAILS;
  let docStore: any;
  let helmReleaseData: Helm.Release;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Helm List", url: routeString.helm },
    { index: 1, name: params.name },
  ];

  $: helmReleaseData = $dataGet;

  $dataSend = [
    {
      type: KubeDataOpType.helmGet,
      request: {
        namespace: params.namespace,
        name: params.name,
      },
    },
  ];
</script>

<HeaderElement>
  <Tabs slot="tabs" bind:activeTab {tabItems} />
  <ResourceToolbar
    slot="toolbar"
    bind:codeMirrorChanged
    bind:toolbarContent
    bind:activeTab
  />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if activeTab === TabIndex.YAML}
    <CodeMirror
      doc={toYaml(helmReleaseData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</div>
