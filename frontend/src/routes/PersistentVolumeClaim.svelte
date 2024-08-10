<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    PersistentVolumeClaimV1GVRK,
    getLabels,
    parseFieldSelector,
    EventV1GVRK,
    randomUUID,
  } from "../lib/util";
  import type { V1PersistentVolumeClaim } from "@kubernetes/client-node";
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
  let docStore: any;
  let persistentVolumeClaimData: V1PersistentVolumeClaim;
  let codeMirrorChanged: boolean;
  let eventListData: any;

  $: toolbarContent = [
    {
      index: 0,
      name: "PersistentVolumeClaim List",
      url: routeString.persistentVolumeClaimList,
    },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: PersistentVolumeClaimV1GVRK,
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
          kind: PersistentVolumeClaimV1GVRK.kind,
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
      persistentVolumeClaimData = dg.data;
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
      persistentVolumeClaimData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: PersistentVolumeClaimV1GVRK,
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
          { name: "Name", value: persistentVolumeClaimData?.metadata?.name },
          {
            name: "Namespace",
            value: persistentVolumeClaimData?.metadata?.namespace,
          },
          {
            name: "Creation Timestamp",
            value: persistentVolumeClaimData?.metadata?.creationTimestamp,
          },
          {
            name: "Finalizers",
            value:
              JSON.stringify(
                persistentVolumeClaimData?.metadata?.finalizers,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
          {
            name: "UID",
            value: persistentVolumeClaimData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(persistentVolumeClaimData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(persistentVolumeClaimData?.metadata?.labels)}
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
                persistentVolumeClaimData?.spec?.accessModes,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Resources",
            value:
              JSON.stringify(
                persistentVolumeClaimData?.spec?.resources,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "storageClassName",
            value: persistentVolumeClaimData?.spec?.storageClassName,
          },
          {
            name: "volumeMode",
            value: persistentVolumeClaimData?.spec?.volumeMode,
          },
          {
            name: "volumeName",
            value: persistentVolumeClaimData?.spec?.volumeName,
          },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Access Modes",
            value: persistentVolumeClaimData?.status?.accessModes,
          },
          {
            name: "Capacity",
            value:
              JSON.stringify(
                persistentVolumeClaimData?.status?.capacity,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Phase",
            value: persistentVolumeClaimData?.status?.phase,
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(persistentVolumeClaimData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
