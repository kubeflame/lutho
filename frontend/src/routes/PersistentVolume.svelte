<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    PersistentVolumeV1GVRK,
    getLabels,
    EventV1GVRK,
    parseFieldSelector,
    randomUUID,
  } from "../lib/util";
  import type { V1PersistentVolume } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import Details from "../lib/Details.svelte";
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
  let store: any;
  let persistentVolumeData: V1PersistentVolume;
  let codeMirrorChanged: boolean;
  let eventListData: any;

  $: toolbarContent = [
    {
      index: 0,
      name: "PersistentVolume List",
      url: routeString.persistentVolumeList,
    },
    { index: 1, name: params.name },
  ];
  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      name: params.name,
      kubeGVRK: PersistentVolumeV1GVRK,
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
          kind: PersistentVolumeV1GVRK.kind,
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
      persistentVolumeData = dg.data;
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
      persistentVolumeData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          name: params.name,
          kubeGVRK: PersistentVolumeV1GVRK,
          data: JSON.stringify(fromYaml($store)),
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
          { name: "Name", value: persistentVolumeData?.metadata?.name },
          {
            name: "Namespace",
            value: persistentVolumeData?.metadata?.namespace,
          },
          {
            name: "Creation Timestamp",
            value: persistentVolumeData?.metadata?.creationTimestamp,
          },
          {
            name: "Finalizers",
            value:
              JSON.stringify(
                persistentVolumeData?.metadata?.finalizers,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
          {
            name: "UID",
            value: persistentVolumeData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(persistentVolumeData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(persistentVolumeData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Access Mode",
            value:
              JSON.stringify(
                persistentVolumeData?.spec?.accessModes,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Capacity",
            value:
              JSON.stringify(persistentVolumeData?.spec?.capacity, null, 2) ??
              "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Claim Ref",
            value:
              JSON.stringify(persistentVolumeData?.spec?.claimRef, null, 2) ??
              "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Host Path",
            value:
              JSON.stringify(persistentVolumeData?.spec?.hostPath, null, 2) ??
              "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Node Affinity",
            value:
              JSON.stringify(
                persistentVolumeData?.spec?.nodeAffinity,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Persistent Volume Reclaim Policy",
            value: persistentVolumeData?.spec?.persistentVolumeReclaimPolicy,
          },
          {
            name: "Storage Class Name",
            value: persistentVolumeData?.spec?.storageClassName,
          },
          {
            name: "Volume Mode",
            value: persistentVolumeData?.spec?.volumeMode,
          },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Phase",
            value: persistentVolumeData?.status?.phase,
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(persistentVolumeData)}
      bind:codeMirrorChanged
      bind:docStore={store}
    />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
