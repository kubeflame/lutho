<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { tabs, routeString, StatefulSetV1GVRK } from "../lib/util";
  import { KubeDataOpType, TabIndex } from "../lib/types";
  import type { V1StatefulSet } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import { fade } from "svelte/transition";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let docStore: any;
  let statefulSetData: V1StatefulSet;
  let activeTab: number = TabIndex.DETAILS;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "StatefulSet List", url: routeString.statefulSetList },
    { index: 1, name: params.name },
  ];

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        namespace: params.namespace,
        name: params.name,
        ...StatefulSetV1GVRK,
      },
    },
  ];

  $: statefulSetData = $dataGet;

  dataUpdate.subscribe((du) => {
    statefulSetData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          namespace: params.namespace,
          name: params.name,
          ...StatefulSetV1GVRK,
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
      doc={toYaml(statefulSetData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</div>
