<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { tabs, routeString, DeploymentV1GVRK } from "../lib/util";
  import { TabIndex, KubeDataOpType } from "../lib/types";
  import type { V1Deployment } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import type { Writable } from "svelte/store";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let activeTab: number = TabIndex.DETAILS;
  let docStore: Writable<string>;
  let deployData: V1Deployment;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Deployment List", url: routeString.deployList },
    { index: 1, name: params.name },
  ];

  $: deployData = $dataGet;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        namespace: params.namespace,
        name: params.name,
        ...DeploymentV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    deployData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          namespace: params.namespace,
          name: params.name,
          ...DeploymentV1GVRK,
          data: JSON.stringify(fromYaml($docStore)),
        },
      },
    ];
  }
</script>

<HeaderElement>
  <Tabs slot="tabs" bind:activeTab {tabItems} />
</HeaderElement>

<RouterPage>
  <ResourceToolbar
    slot="resource-toolbar"
    bind:codeMirrorChanged
    bind:toolbarContent
    bind:activeTab
    onClickSubmit={update}
  />
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if activeTab === TabIndex.YAML}
    <CodeMirror doc={toYaml(deployData)} bind:docStore bind:codeMirrorChanged />
  {/if}
</RouterPage>
