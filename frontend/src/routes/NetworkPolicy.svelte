<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    NetworkPolicyV1GVRK,
    getLabels,
  } from "../lib/util";
  import type { V1NetworkPolicy } from "@kubernetes/client-node";
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
  let networkPolicyData: V1NetworkPolicy;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    {
      index: 0,
      name: "NetworkPolicy List",
      url: routeString.networkPolicyList,
    },
    { index: 1, name: params.name },
  ];

  $: networkPolicyData = $dataGet;

  $dataSend = [
    {
      type: "get",
      request: {
        namespace: params.namespace,
        name: params.name,
        kubeGVRK: NetworkPolicyV1GVRK,
      },
    },
  ];

  dataUpdate.subscribe((du) => {
    networkPolicyData = du;
  });

  function update() {
    $dataSend = [
      {
        type: "update",
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: NetworkPolicyV1GVRK,
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
          { name: "Name", value: networkPolicyData?.metadata?.name },
          { name: "Namespace", value: networkPolicyData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: networkPolicyData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: networkPolicyData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(networkPolicyData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(networkPolicyData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Egress",
            value:
              JSON.stringify(networkPolicyData?.spec?.egress, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Ingress",
            value:
              JSON.stringify(networkPolicyData?.spec?.ingress, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Pod Selector",
            value:
              JSON.stringify(networkPolicyData?.spec?.podSelector, null, 2) ??
              "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Policy Types",
            value:
              JSON.stringify(networkPolicyData?.spec?.policyTypes, null, 2) ??
              "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(networkPolicyData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
