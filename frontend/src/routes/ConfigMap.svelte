<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import { tabs, routeString, ConfigMapV1GVRK, getLabels } from "../lib/util";
  import type { V1ConfigMap } from "@kubernetes/client-node";
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
  let configMapData: V1ConfigMap;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "ConfigMap List", url: routeString.configMapList },
    { index: 1, name: params.name },
  ];

  $: configMapData = $dataGet;

  $dataSend = [
    {
      type: "get",
      request: {
        namespace: params.namespace,
        name: params.name,
        kubeGVRK: ConfigMapV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    configMapData = du;
  });

  function update() {
    $dataSend = [
      {
        type: "update",
        request: {
          name: params.name,
          namespace: params.namespace,
          kubeGVRK: ConfigMapV1GVRK,
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
          { name: "Name", value: configMapData?.metadata?.name },
          { name: "Namespace", value: configMapData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: configMapData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: configMapData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(configMapData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(configMapData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Data"}>
      <div
        class="border-base-300 bg-base-200/10 overflow-hidden break-words
          rounded-lg border px-1 font-mono text-sm"
      >
        {#if configMapData?.data}
          {#each Object.entries(configMapData.data) as data}
            <div>
              {data[0]}:
              <div class="ml-4 whitespace-pre-wrap">{data[1]}</div>
            </div>
          {/each}
        {/if}
      </div>
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(configMapData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
