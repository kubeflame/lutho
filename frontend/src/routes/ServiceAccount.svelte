<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    ServiceAccountV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1ServiceAccount } from "@kubernetes/client-node";
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

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let serviceAccountData: V1ServiceAccount;
  let codeMirrorChanged: boolean;

  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  $: toolbarContent = [
    {
      index: 0,
      name: "ServiceAccount List",
      url: routeString.serviceAccountList,
    },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: ServiceAccountV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      serviceAccountData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      serviceAccountData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: ServiceAccountV1GVRK,
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
          { name: "Name", value: serviceAccountData?.metadata?.name },
          { name: "Namespace", value: serviceAccountData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: serviceAccountData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: serviceAccountData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(serviceAccountData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(serviceAccountData?.metadata?.labels)}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(serviceAccountData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
