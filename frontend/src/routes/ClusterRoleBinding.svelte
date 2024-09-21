<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    ClusterRoleBindingV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1ClusterRoleBinding } from "@kubernetes/client-node";
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
  let clusterRoleBindingData: V1ClusterRoleBinding;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    {
      index: 0,
      name: "ClusterRoleBinding List",
      url: routeString.clusterRoleBindingList,
    },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      name: params.name,
      kubeGVRK: ClusterRoleBindingV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      clusterRoleBindingData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      clusterRoleBindingData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          name: params.name,
          kubeGVRK: ClusterRoleBindingV1GVRK,
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
          { name: "Name", value: clusterRoleBindingData?.metadata?.name },
          {
            name: "Creation Timestamp",
            value: clusterRoleBindingData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: clusterRoleBindingData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(clusterRoleBindingData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(clusterRoleBindingData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Role Reference"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Name",
            value: clusterRoleBindingData?.roleRef?.name,
          },
          {
            name: "API Group",
            value: clusterRoleBindingData?.roleRef?.apiGroup,
          },
          {
            name: "Kind",
            value: clusterRoleBindingData?.roleRef?.kind,
          },
        ]}
      />
    </Details>

    <Details title={"Subjects"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[
          { name: "Name" },
          { name: "Namespace" },
          { name: "Kind" },
          { name: "API Group" },
        ]}
      >
        <tbody class="">
          {#if clusterRoleBindingData?.subjects}
            {#each clusterRoleBindingData?.subjects as subject}
              <tr>
                <td>{subject.name ?? "-"}</td>
                <td>{subject.namespace ?? "-"}</td>
                <td>{subject.kind ?? "-"}</td>
                <td>{subject.apiGroup ?? "-"}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(clusterRoleBindingData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
