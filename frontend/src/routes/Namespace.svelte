<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { tabs, routeString, NamespaceV1GVRK } from "../lib/util";
  import { TabIndex, KubeDataOpType } from "../lib/types";
  import type { V1Namespace } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let activeTab: number = TabIndex.DETAILS;
  let docStore: any;
  let namespaceData: V1Namespace;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Namespace List", url: routeString.namespaceList },
    { index: 1, name: params.name },
  ];

  $: namespaceData = $dataGet;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        name: params.name,
        ...NamespaceV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    namespaceData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          name: params.name,
          ...NamespaceV1GVRK,
          data: JSON.stringify(fromYaml($docStore)),
        },
      },
    ];
  }
</script>

<HeaderElement>
  <Tabs slot="tabs" bind:activeTab {tabItems} />
  <ResourceToolbar
    slot="toolbar"
    bind:codeMirrorChanged
    bind:toolbarContent
    bind:activeTab
    onClickSubmit={update}
  />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if activeTab === TabIndex.YAML}
    <CodeMirror
      doc={toYaml(namespaceData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</div>
