<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import { tabs, routeString, EventV1GVRK, NodeV1GVRK } from "../lib/util";
  import { TabIndex, KubeDataOpType } from "../lib/types";
  import Events from "../lib/Events.svelte";
  import type { EventsV1EventList, V1Node } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import { fade } from "svelte/transition";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import Details from "../lib/Details.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  let nodeEventOpts = {
    fieldSelector: `regarding.kind=Node,regarding.name=${params.name}`,
  };

  const tabItems = [tabs.details, tabs.yaml, tabs.events];
  const { sockError, isLoading, dataSend, dataList, dataGet, dataUpdate } =
    socketStore();

  let docStore: any;
  let nodeData: V1Node;
  let activeTab: number = TabIndex.DETAILS;
  let eventListData: EventsV1EventList;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Node List", url: routeString.nodeList },
    { index: 1, name: params.name },
  ];

  $: nodeData = $dataGet;
  $: eventListData = $dataList;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        name: params.name,
        ...NodeV1GVRK,
      },
    },
    {
      type: KubeDataOpType.list,
      request: {
        name: params.name,
        ...EventV1GVRK,
        options: nodeEventOpts,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    nodeData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          name: params.name,
          ...NodeV1GVRK,
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
  {:else if activeTab === TabIndex.DETAILS}
    <!-- <Details {params} /> -->
  {:else if activeTab === TabIndex.YAML}
    <CodeMirror doc={toYaml(nodeData)} bind:codeMirrorChanged bind:docStore />
  {:else if activeTab === TabIndex.EVENTS}
    <Events kubeEvents={eventListData.items} />
  {/if}
</div>
