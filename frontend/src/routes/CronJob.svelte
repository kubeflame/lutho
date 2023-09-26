<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    CronJobV1GVRK,
    getLabels,
    EventV1GVRK,
    parseFieldSelector,
  } from "../lib/util";
  import type { V1CronJob } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Logs from "../lib/Logs.svelte";
  import Details from "../lib/Details.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";
  import Events from "../lib/Events.svelte";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml, tabs.events];
  const {
    sockState,
    sockError,
    isLoading,
    dataSend,
    dataList,
    dataGet,
    dataUpdate,
  } = socketStore();

  let tabQueryParam: TabQueryParam;
  let followLogs = false;
  let stopLogs = false;
  let refreshLogsComponent = false;
  let docStore: any;
  let cronJobData: V1CronJob;
  let codeMirrorChanged: boolean;
  let containers: string[];

  $: toolbarContent = [
    { index: 0, name: "CronJob List", url: routeString.cronJobList },
    { index: 1, name: params.name },
  ];
  $: cronJobData = $dataGet;
  $: eventListData = $dataList;

  $dataSend = [
    {
      type: "get",
      request: {
        namespace: params.namespace,
        name: params.name,
        kubeGVRK: CronJobV1GVRK,
      },
    },
  ];

  $: if (tabQueryParam === "events") {
    $dataSend = [
      {
        type: "list",
        request: {
          name: params.name,
          kubeGVRK: EventV1GVRK,
          kubeOptions: {
            fieldSelector: parseFieldSelector({
              kind: CronJobV1GVRK.kind,
              name: params.name,
              namespace: params.namespace,
            }),
          },
        },
      },
    ];
  }

  dataSend.subscribe((ds) => {
    if (ds && $sockState.state === WebSocket.CLOSED) $sockState.refresh = true;
  });

  dataUpdate.subscribe((du) => {
    cronJobData = du;
  });

  function update() {
    $dataSend = [
      {
        type: "update",
        request: {
          name: params.name,
          namespace: params.namespace,
          kubeGVRK: CronJobV1GVRK,
          data: JSON.stringify(fromYaml($docStore)),
        },
      },
    ];
  }
</script>

<HeaderElement>
  <Tabs
    slot="tabs"
    bind:tabQueryParam
    bind:containers
    tabQueryParamDefault={"details"}
    {tabItems}
  />
</HeaderElement>

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar
    slot="resource-toolbar"
    bind:codeMirrorChanged
    bind:tabQueryParam
    onClickSubmit={update}
  >
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  {#if tabQueryParam === "details"}
    <Details title={"Resource Information"}>
      <EmbeddedTable
        tagName={"Summary"}
        tableType={"custom-vertical"}
        tableItems={[
          { name: "Name", value: cronJobData?.metadata?.name },
          { name: "Namespace", value: cronJobData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: cronJobData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: cronJobData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(cronJobData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(cronJobData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Schedule",
            value: JSON.stringify(cronJobData?.spec?.schedule, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Successful Jobs History Limit",
            value: cronJobData?.spec?.successfulJobsHistoryLimit,
          },
          {
            name: "Failed Jobs History Limit",
            value: cronJobData?.spec?.failedJobsHistoryLimit,
          },
          {
            name: "Suspend",
            value: cronJobData?.spec?.suspend,
          },
          {
            name: "Concurrency Policy",
            value: cronJobData?.spec?.concurrencyPolicy,
          },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Last Schedule Time",
            value: cronJobData?.status?.lastScheduleTime,
          },
          {
            name: "Last Successful Time",
            value: cronJobData?.status?.lastSuccessfulTime,
          },
        ]}
      />
    </Details>

    <Details title={"Template"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Name",
            value: cronJobData?.spec?.jobTemplate?.metadata?.name,
          },
          {
            name: "Namespace",
            value: cronJobData?.spec?.jobTemplate?.metadata?.namespace,
          },
          {
            name: "Creation Timestamp",
            value: cronJobData?.spec?.jobTemplate?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: cronJobData?.spec?.jobTemplate?.metadata?.uid,
          },
        ]}
      />

      {#if cronJobData?.spec?.jobTemplate?.spec?.template?.spec?.containers}
        {#each cronJobData?.spec?.jobTemplate?.spec?.template?.spec?.containers as container}
          <EmbeddedTable
            tagName={"Containers"}
            tableType={"custom-vertical"}
            tableItems={[
              {
                name: "Name",
                value: container.name,
              },
              {
                name: "Image",
                value: container.image,
              },
              {
                name: "Image Pull Policy",
                value: container.imagePullPolicy,
              },
              {
                name: "Command",
                value: JSON.stringify(container.command, null, 2) ?? "-",
                className:
                  "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
              },
              {
                name: "Resources",
                value: JSON.stringify(container.resources, null, 2) ?? "-",
                className:
                  "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
              },
            ]}
          />
        {/each}
      {/if}

      <EmbeddedTable
        tagName={"Other info"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "DNS Policy",
            value:
              cronJobData?.spec?.jobTemplate?.spec?.template?.spec?.dnsPolicy,
          },
          {
            name: "Restart Policy",
            value:
              cronJobData?.spec?.jobTemplate?.spec?.template?.spec
                ?.restartPolicy,
          },
          {
            name: "Scheduler Name",
            value:
              cronJobData?.spec?.jobTemplate?.spec?.template?.spec
                ?.schedulerName,
          },
          {
            name: "Security Context",
            value:
              JSON.stringify(
                cronJobData?.spec?.jobTemplate?.spec?.template?.spec
                  ?.securityContext,
                null,
                2,
              ) ?? "-",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(cronJobData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
