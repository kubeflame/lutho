<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { tabs, routeString, CronJobV1GVRK } from "../lib/util";
  import { KubeDataOpType, TabIndex } from "../lib/types";
  import type { V1CronJob } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Logs from "../lib/Logs.svelte";
  import Details from "../lib/Details.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataList, dataGet, dataUpdate } =
    socketStore();

  let activeTab: number = TabIndex.DETAILS;
  let followLogs = false;
  let stopLogs = false;
  let refreshLogsComponent = false;
  let docStore: any;
  let cronJobData: V1CronJob;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "CronJob List", url: routeString.cronJobList },
    { index: 1, name: params.name },
  ];

  $: cronJobData = $dataGet;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        namespace: params.namespace,
        name: params.name,
        ...CronJobV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    cronJobData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          name: params.name,
          namespace: params.namespace,
          ...CronJobV1GVRK,
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
      doc={toYaml(cronJobData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</div>
