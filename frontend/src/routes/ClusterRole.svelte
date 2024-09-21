<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    ClusterRoleV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1ClusterRole } from "@kubernetes/client-node";
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
  let clusterRoleData: V1ClusterRole;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    {
      index: 0,
      name: "ClusterRole List",
      url: routeString.clusterRoleList,
    },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      name: params.name,
      kubeGVRK: ClusterRoleV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      clusterRoleData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      clusterRoleData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          name: params.name,
          kubeGVRK: ClusterRoleV1GVRK,
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
          { name: "Name", value: clusterRoleData?.metadata?.name },
          {
            name: "Creation Timestamp",
            value: clusterRoleData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: clusterRoleData?.metadata?.uid,
          },
        ]}
      />

      {#if clusterRoleData?.aggregationRule}
        <EmbeddedTable
          tagName={"Aggregation Rule"}
          tableType={"custom-vertical"}
          tableItems={[
            {
              name: "ClusterRole Selectors",
              value:
                JSON.stringify(
                  clusterRoleData?.aggregationRule.clusterRoleSelectors,
                  null,
                  2,
                ) ?? "-",
              className:
                "max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
            },
          ]}
        />
      {/if}

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(clusterRoleData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(clusterRoleData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Rules"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[
          { name: "API Groups" },
          { name: "Resources" },
          { name: "Verbs" },
        ]}
      >
        <tbody class="text-sm">
          {#if clusterRoleData?.rules}
            {#each clusterRoleData.rules as rule}
              <tr>
                <td
                  class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
                >
                  {JSON.stringify(rule.apiGroups, null, 2) ?? "-"}</td
                >
                <td
                  class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
                >
                  {JSON.stringify(rule.resources, null, 2) ?? "-"}</td
                >
                <td
                  class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
                >
                  {JSON.stringify(rule.verbs, null, 2) ?? "-"}</td
                >
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(clusterRoleData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
