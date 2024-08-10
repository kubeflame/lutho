<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    RoleV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1Role } from "@kubernetes/client-node";
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
  const { sockError, isLoading, dataSend, dataList, dataGet, dataUpdate } =
    socketStore();

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let roleData: V1Role;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Role List", url: routeString.roleList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: RoleV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      roleData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      roleData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: RoleV1GVRK,
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
          { name: "Name", value: roleData?.metadata?.name },
          {
            name: "Creation Timestamp",
            value: roleData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: roleData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(roleData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(roleData?.metadata?.labels)}
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
          {#if roleData?.rules}
            {#each roleData.rules as rule}
              <tr>
                <td
                  class="border-base-300 bg-base-200/10 max-h-fit max-w-fit rounded-lg border px-1 font-mono text-sm"
                >
                  {JSON.stringify(rule.apiGroups, null, 2) ?? "-"}</td
                >
                <td
                  class="border-base-300 bg-base-200/10 max-h-fit max-w-fit rounded-lg border px-1 font-mono text-sm"
                >
                  {JSON.stringify(rule.resources, null, 2) ?? "-"}</td
                >
                <td
                  class="border-base-300 bg-base-200/10 max-h-fit max-w-fit rounded-lg border px-1 font-mono text-sm"
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
    <CodeMirror value={toYaml(roleData)} bind:codeMirrorChanged bind:docStore />
  {/if}
</RouterPage>
