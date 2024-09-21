<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    RoleBindingV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1RoleBinding } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";
  import Details from "../lib/Details.svelte";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let roleBindingData: V1RoleBinding;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "RoleBinding List", url: routeString.roleBindingList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: RoleBindingV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      roleBindingData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      roleBindingData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: RoleBindingV1GVRK,
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
          { name: "Name", value: roleBindingData?.metadata?.name },
          {
            name: "Creation Timestamp",
            value: roleBindingData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: roleBindingData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(roleBindingData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(roleBindingData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Role Ref"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[{ name: "Name", value: roleBindingData?.roleRef?.name }]}
      />
    </Details>

    <Details title={"Subjects"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[{ name: "Kind" }, { name: "Name" }, { name: "Namespace" }]}
      >
        <tbody class="text-sm">
          {#if roleBindingData?.subjects}
            {#each roleBindingData.subjects as subject}
              <tr>
                <td
                  class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
                >
                  {subject.kind}
                </td>
                <td
                  class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
                >
                  {subject.name}
                </td>
                <td
                  class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
                >
                  {subject.namespace}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(roleBindingData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
