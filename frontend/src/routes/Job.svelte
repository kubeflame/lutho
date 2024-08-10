<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    JobV1GVRK,
    getLabels,
    EventV1GVRK,
    parseFieldSelector,
    randomUUID,
  } from "../lib/util";
  import type { V1Job } from "@kubernetes/client-node";
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
  import CaptionTag from "../lib/CaptionTag.svelte";
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
  let jobData: V1Job;
  let codeMirrorChanged: boolean;
  let containers: string[];
  let eventListData: any;

  $: toolbarContent = [
    { index: 0, name: "Job List", url: routeString.jobList },
    { index: 1, name: params.name },
  ];
  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: JobV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      name: params.name,
      kubeGVRK: EventV1GVRK,
      kubeOptions: {
        fieldSelector: parseFieldSelector({
          kind: JobV1GVRK.kind,
          name: params.name,
          namespace: params.namespace,
        }),
      },
    },
  } as any;

  $: if (tabQueryParam === "events") {
    $dataSend = [sendList];
  }

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      jobData = dg.data;
    }
  });

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      eventListData = dl.data;
    }
  });

  dataSend.subscribe((ds) => {
    if (ds && $sockState.state === WebSocket.CLOSED) $sockState.refresh = true;
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      jobData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: JobV1GVRK,
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
          { name: "Name", value: jobData?.metadata?.name },
          { name: "Namespace", value: jobData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: jobData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: jobData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(jobData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(jobData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Backoff Limit",
            value: jobData?.spec?.backoffLimit,
          },
          {
            name: "Completion Mode",
            value: jobData?.spec?.completionMode,
          },
          {
            name: "Completions",
            value: jobData?.spec?.completions,
          },
          {
            name: "Parallelism",
            value: jobData?.spec?.parallelism,
          },
          {
            name: "Selector",
            value: JSON.stringify(jobData?.spec?.selector, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Suspend",
            value: jobData?.spec?.suspend,
          },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Completion Time",
            value: jobData?.status?.completionTime,
          },
          {
            name: "Conditions",
            value: JSON.stringify(jobData?.status?.conditions, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Ready",
            value: jobData?.status?.ready,
          },
          {
            name: "Start Time",
            value: jobData?.status?.startTime,
          },
          {
            name: "Succeeded",
            value: jobData?.status?.succeeded,
          },
          {
            name: "Uncounted Terminated Pods",
            value:
              JSON.stringify(
                jobData?.status?.uncountedTerminatedPods,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
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
            value: jobData?.spec?.template?.metadata?.name,
          },
          {
            name: "Namespace",
            value: jobData?.spec?.template?.metadata?.namespace,
          },
          {
            name: "Creation Timestamp",
            value: jobData?.spec?.template?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: jobData?.spec?.template?.metadata?.uid,
          },
        ]}
      />

      {#if jobData?.spec?.template?.spec?.containers}
        <div class="card card-compact mt-2 rounded-lg">
          <CaptionTag tagName={"Containers"} />
          {#each jobData?.spec?.template?.spec?.containers as container}
            <div class="card-title pl-2 pt-2 text-sm font-light tracking-wider">
              <div class="badge badge-ghost badge-xs" />
              {container.name}
            </div>

            <div class="card-body space-y-3">
              <EmbeddedTable
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
                    name: "Args",
                    value: JSON.stringify(container.args, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                  {
                    name: "Env",
                    value: JSON.stringify(container.env, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                  {
                    name: "Image Pull Policy",
                    value: container.imagePullPolicy,
                  },
                  {
                    name: "Volume Mounts",
                    value:
                      JSON.stringify(container.volumeMounts, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                ]}
              />
            </div>
          {/each}
        </div>
      {/if}

      <EmbeddedTable
        tagName={"Other info"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "DNS Policy",
            value: jobData?.spec?.template?.spec?.dnsPolicy,
          },
          {
            name: "Node Selector",
            value:
              JSON.stringify(
                jobData?.spec?.template?.spec?.nodeSelector,
                null,
                2,
              ) ?? "-",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Restart Policy",
            value: jobData?.spec?.template?.spec?.restartPolicy,
          },
          {
            name: "Scheduler Name",
            value: jobData?.spec?.template?.spec?.schedulerName,
          },
          {
            name: "Security Context",
            value:
              JSON.stringify(
                jobData?.spec?.template?.spec?.securityContext,
                null,
                2,
              ) ?? "-",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Service Account",
            value: jobData?.spec?.template?.spec?.serviceAccount,
          },
          {
            name: "Volumes",
            value:
              JSON.stringify(jobData?.spec?.template?.spec?.volumes, null, 2) ??
              "-",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror value={toYaml(jobData)} bind:codeMirrorChanged bind:docStore />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
