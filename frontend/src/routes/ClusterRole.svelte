<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { tabs, routeString, ClusterRoleV1GVRK } from "../lib/util";
  import { TabIndex, KubeDataOpType } from "../lib/types";
  import type { V1ClusterRole } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let activeTab: number = TabIndex.DETAILS;
  let docStore: any;
  let clusterRoleData: V1ClusterRole;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    {
      index: 0,
      name: "ClusterRole List",
      url: routeString.clusterRoleList,
    },
    { index: 1, name: params.name },
  ];

  $: clusterRoleData = $dataGet;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        name: params.name,
        ...ClusterRoleV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    clusterRoleData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          name: params.name,
          ...ClusterRoleV1GVRK,
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
      doc={toYaml(clusterRoleData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</div>
