<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { tabs, routeString, ConfigMapV1GVRK } from "../lib/util";
  import { TabIndex, KubeDataOpType } from "../lib/types";
  import type { V1ConfigMap } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import { fade } from "svelte/transition";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let activeTab: number = TabIndex.DETAILS;
  let docStore: any;
  let configMapData: V1ConfigMap;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "ConfigMap List", url: routeString.configMapList },
    { index: 1, name: params.name },
  ];

  $: configMapData = $dataGet;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        namespace: params.namespace,
        name: params.name,
        ...ConfigMapV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    configMapData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          name: params.name,
          namespace: params.namespace,
          ...ConfigMapV1GVRK,
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
  {#if activeTab === TabIndex.YAML}
    {#if $sockError}
      <ErrorPage bind:errorMessage={$sockError} />
    {:else if $isLoading}
      <LoadingNewton />
    {:else}
      <CodeMirror
        doc={toYaml(configMapData)}
        bind:codeMirrorChanged
        bind:docStore
      />
    {/if}
  {/if}
</div>
