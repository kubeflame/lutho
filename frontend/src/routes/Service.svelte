<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    ServiceV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1Service } from "@kubernetes/client-node";
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
  let serviceData: V1Service;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Service List", url: routeString.serviceList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: ServiceV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      serviceData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      serviceData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: ServiceV1GVRK,
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
          { name: "Name", value: serviceData?.metadata?.name },
          { name: "Namespace", value: serviceData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: serviceData?.metadata?.creationTimestamp,
          },
          {
            name: "Finalizerz",
            value:
              JSON.stringify(serviceData?.metadata?.finalizers, null, 2) ??
              "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
          {
            name: "UID",
            value: serviceData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(serviceData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(serviceData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Type",
            value: serviceData?.spec?.type,
          },
          {
            name: "Allocate Load Balancer Node Ports",
            value: serviceData?.spec?.allocateLoadBalancerNodePorts,
          },
          {
            name: "External Traffic Policy",
            value: serviceData?.spec?.externalTrafficPolicy,
          },
          {
            name: "Internal Traffic Policy",
            value: serviceData?.spec?.internalTrafficPolicy,
          },
          {
            name: "Session Affinity",
            value: serviceData?.spec?.sessionAffinity,
          },
          {
            name: "Cluster IPs",
            value:
              JSON.stringify(serviceData?.spec?.clusterIPs, null, 2) ?? "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
          {
            name: "Ports",
            value: JSON.stringify(serviceData?.spec?.ports, null, 2) ?? "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
          {
            name: "Selector",
            value: JSON.stringify(serviceData?.spec?.selector, null, 2) ?? "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Load Balancer",
            value:
              JSON.stringify(serviceData?.status?.loadBalancer, null, 2) ??
              "{}",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(serviceData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
