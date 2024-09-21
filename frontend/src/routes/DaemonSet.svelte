<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    DaemonSetV1GVRK,
    getLabels,
    EventV1GVRK,
    parseFieldSelector,
    randomUUID,
  } from "../lib/util";
  import type { V1DaemonSet } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import Details from "../lib/Details.svelte";
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
  let docStore: any;
  let daemonSetData: V1DaemonSet;
  let codeMirrorChanged: boolean;
  let eventListData: any;

  $: toolbarContent = [
    { index: 0, name: "DaemonSet List", url: routeString.daemonSetList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: DaemonSetV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      name: params.name,
      namespace: params.namespace,
      kubeGVRK: EventV1GVRK,
      kubeOptions: {
        fieldSelector: parseFieldSelector({
          kind: DaemonSetV1GVRK.kind,
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
      daemonSetData = dg.data;
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
      daemonSetData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: DaemonSetV1GVRK,
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
    tabQueryParamDefault={"details"}
    {tabItems}
  />
</HeaderElement>

<RouterPage bind:errorMessage={$sockError} bind:loading={$isLoading}>
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
          { name: "Name", value: daemonSetData?.metadata?.name },
          { name: "Namespace", value: daemonSetData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: daemonSetData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: daemonSetData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tagName={"Status"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Current Number Scheduled",
            value: daemonSetData?.status?.currentNumberScheduled,
          },
          {
            name: "Desired Number Scheduled",
            value: daemonSetData?.status?.desiredNumberScheduled,
          },
          {
            name: "Number Available",
            value: daemonSetData?.status?.numberAvailable,
          },
          {
            name: "Number Misscheduled",
            value: daemonSetData?.status?.numberMisscheduled,
          },
          {
            name: "Observed Generation",
            value: daemonSetData?.status?.observedGeneration,
          },
          {
            name: "Updated Number Scheduled",
            value: daemonSetData?.status?.updatedNumberScheduled,
          },
        ]}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Revision History Limit",
            value: daemonSetData?.spec?.revisionHistoryLimit,
          },
          {
            name: "Selector",
            value:
              JSON.stringify(daemonSetData?.spec?.selector, null, 2) ?? "-",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
          {
            name: "Update Strategy",
            value:
              JSON.stringify(daemonSetData?.spec?.updateStrategy, null, 2) ??
              "-",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
        ]}
      />
    </Details>

    <Details title={"Template"}>
      <EmbeddedTable
        tagName={"Details"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "DNS Policy",
            value: daemonSetData?.spec?.template?.spec?.dnsPolicy,
          },
          {
            name: "Restart Policy",
            value: daemonSetData?.spec?.template?.spec?.restartPolicy,
          },
          {
            name: "Scheduler Name",
            value: daemonSetData?.spec?.template?.spec?.schedulerName,
          },
          {
            name: "Security Context",
            value:
              JSON.stringify(
                daemonSetData?.spec?.template?.spec?.securityContext,
                null,
                2,
              ) ?? "-",
            className:
              "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Service Account",
            value: daemonSetData?.spec?.template?.spec?.serviceAccount,
          },
          {
            name: "Service Account Name",
            value: daemonSetData?.spec?.template?.spec?.serviceAccountName,
          },
          {
            name: "Tolerations",
            value:
              JSON.stringify(
                daemonSetData?.spec?.template?.spec?.tolerations,
                null,
                2,
              ) ?? "-",
            className:
              "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
        ]}
      />

      {#if daemonSetData?.spec?.template?.spec?.containers}
        <div class="card card-compact mt-2 rounded-lg">
          <CaptionTag tagName={"Containers"} />
          {#each daemonSetData?.spec?.template?.spec?.containers as container}
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
                    name: "Image Pull Policy",
                    value: container.imagePullPolicy,
                  },
                  {
                    name: "Ports",
                    value: JSON.stringify(container.ports, null, 2) ?? "-",
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
                    name: "Security Context",
                    value:
                      JSON.stringify(container.securityContext, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                ]}
              />
            </div>
          {/each}
        </div>
      {/if}
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(daemonSetData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
