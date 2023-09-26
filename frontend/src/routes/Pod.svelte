<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import { embeddedActiveTab } from "../lib/stores";
  import Shell from "../lib/Shell.svelte";
  import { routeString, tabs, EventV1GVRK, PodV1GVRK } from "../lib/util";
  import { KubeDataOpType, TabIndex } from "../lib/types";
  import Events from "../lib/Events.svelte";
  import Details from "../lib/Details.svelte";
  import type {
    EventsV1EventList,
    V1Container,
    V1Pod,
  } from "@kubernetes/client-node";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { onMount } from "svelte";
  import Logs from "../lib/Logs.svelte";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import { fade } from "svelte/transition";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";

  export let params: any;

  let docStore: any;
  let shellEmbed: boolean;
  let followLogs: boolean = false;
  let stopLogs: boolean = false;
  let refreshLogsComponent: boolean = false;
  let refreshShellComponent: boolean = false;
  let podData: V1Pod;
  let eventListData: EventsV1EventList;
  let activeTab: number = TabIndex.DETAILS;
  let activeContainer: string = "";
  let showShellReconnect: boolean;
  let codeMirrorChanged: boolean;
  let eventsReversed: number = 0;

  const tabItems = [
    tabs.details,
    tabs.yaml,
    tabs.logs,
    tabs.shell,
    tabs.events,
  ];
  const podEventOpts = {
    fieldSelector: `regarding.kind=Pod,regarding.name=${params.name},regarding.namespace=${params.namespace}`,
  };
  const getContainers = (containers: V1Container[]): string[] => {
    return containers.map((c) => {
      return c.name;
    });
  };

  const { sockError, isLoading, dataSend, dataList, dataGet, dataUpdate } =
    socketStore();

  onMount(async () => {
    if ($embeddedActiveTab === TabIndex.SHELL_EMBED) {
      activeTab = TabIndex.SHELL;
      shellEmbed = true;
      $embeddedActiveTab = undefined;
    } else if ($embeddedActiveTab === TabIndex.LOGS_EMBED) {
      activeTab = TabIndex.LOGS;
    }
  });

  $: toolbarContent = [
    { index: 0, name: "Pod List", url: routeString.podList },
    { index: 1, name: params.name },
  ];
  $: podData = $dataGet;
  $: eventListData = $dataList;

  $dataSend = [
    {
      type: KubeDataOpType.get,
      request: {
        namespace: params.namespace,
        name: params.name,
        ...PodV1GVRK,
      },
    },
    {
      type: KubeDataOpType.list,
      request: {
        namespace: params.namespace,
        name: params.name,
        ...EventV1GVRK,
        options: podEventOpts,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    podData = du;
  });

  function update() {
    $dataSend = [
      {
        type: KubeDataOpType.update,
        request: {
          namespace: params.namespace,
          name: params.name,
          ...PodV1GVRK,
          data: JSON.stringify(fromYaml($docStore)),
        },
      },
    ];
  }

  let unique = {};
  function restart() {
    unique = {};
    refreshShellComponent = false;
  }

  $: if (refreshShellComponent) restart();
</script>

<HeaderElement>
  <Tabs
    slot="tabs"
    bind:activeTab
    bind:refreshLogsComponent
    bind:refreshShellComponent
    bind:activeContainer
    {tabItems}
  />
  <ResourceToolbar
    slot="toolbar"
    bind:codeMirrorChanged
    bind:toolbarContent
    bind:activeTab
    bind:followLogs
    bind:stopLogs
    bind:showShellReconnect
    bind:refreshShellComponent
    onClickSubmit={update}
  />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if activeTab === TabIndex.YAML}
    <CodeMirror doc={toYaml(podData)} bind:codeMirrorChanged bind:docStore />
  {:else if activeTab === TabIndex.LOGS}
    <Logs
      bind:followLogs
      bind:stopLogs
      bind:refreshLogsComponent
      bind:activeContainer
      termParams={{
        name: params.name,
        namespace: params.namespace,
        containers: getContainers(podData?.spec?.containers ?? []),
      }}
    />
  {:else if activeTab === TabIndex.SHELL}
    {#key unique}
      <Shell
        termParams={{
          name: params.name,
          namespace: params.namespace,
          containers: getContainers(podData?.spec?.containers ?? []),
          shellEmbed,
        }}
        bind:showShellReconnect
        bind:activeContainer
      />
    {/key}
  {:else if activeTab === TabIndex.EVENTS}
    <Events bind:kubeEvents={eventListData.items} bind:eventsReversed />
  {/if}
</div>
