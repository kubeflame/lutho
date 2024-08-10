<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    StorageClassV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1StorageClass } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import Details from "../lib/Details.svelte";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let storageClassData: V1StorageClass;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "StorageClass List", url: routeString.storageClassList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      name: params.name,
      kubeGVRK: StorageClassV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      storageClassData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      storageClassData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          name: params.name,
          kubeGVRK: StorageClassV1GVRK,
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
          { name: "Name", value: storageClassData?.metadata?.name },
          { name: "Namespace", value: storageClassData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: storageClassData?.metadata?.creationTimestamp,
          },
          {
            name: "Provisioner",
            value: storageClassData?.provisioner,
          },
          {
            name: "Reclaim Policy",
            value: storageClassData?.reclaimPolicy,
          },
          {
            name: "Volume Binding Mode",
            value: storageClassData?.volumeBindingMode,
          },
          {
            name: "UID",
            value: storageClassData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(storageClassData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(storageClassData?.metadata?.labels)}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(storageClassData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
