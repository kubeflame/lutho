<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    NamespaceV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1Namespace } from "@kubernetes/client-node";
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
  let namespaceData: V1Namespace;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Namespace List", url: routeString.namespaceList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      name: params.name,
      kubeGVRK: NamespaceV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      namespaceData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      namespaceData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          name: params.name,
          kubeGVRK: NamespaceV1GVRK,
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
          { name: "Name", value: namespaceData?.metadata?.name },
          {
            name: "Creation Timestamp",
            value: namespaceData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: namespaceData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(namespaceData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(namespaceData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Finalizers",
            value:
              JSON.stringify(namespaceData?.spec?.finalizers, null, 2) ?? "-",
            className:
              "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
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
            value: namespaceData?.status?.phase,
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(namespaceData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
