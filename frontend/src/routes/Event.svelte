<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    EventV1GVRK,
    jsonStringClassName,
    randomUUID,
  } from "../lib/util";
  import type { CoreV1Event } from "@kubernetes/client-node";
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
  let eventData: CoreV1Event;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    { index: 0, name: "Event List", url: routeString.eventList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: EventV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      eventData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      eventData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: EventV1GVRK,
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
          { name: "Name", value: eventData?.metadata?.name },
          { name: "Namespace", value: eventData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: eventData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: eventData?.metadata?.uid,
          },
        ]}
      />
    </Details>

    <Details title={"Details"}>
      <EmbeddedTable
        tagName={"Involved Object"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Kind",
            value: eventData?.involvedObject.kind,
          },
          {
            name: "Name",
            value: eventData?.involvedObject.name,
          },
          {
            name: "Namespace",
            value: eventData?.involvedObject.namespace,
          },
        ]}
      />

      <EmbeddedTable
        tagName={"Data"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Count",
            value: eventData?.count,
          },
          {
            name: "Kind",
            value: eventData?.kind,
          },
          {
            name: "Type",
            value: eventData?.type,
          },
          {
            name: "Reason",
            value: eventData?.reason,
          },
          {
            name: "Message",
            value: eventData?.message,
            className: "whitespace-pre-wrap",
          },
          {
            name: "Reporting Component",
            value: eventData?.reportingComponent,
          },
          {
            name: "Reporting Instance",
            value: eventData?.reportingInstance,
          },
          {
            name: "Source",
            value: JSON.stringify(eventData?.source, null, 2) ?? "{}",
            className: jsonStringClassName,
          },
        ]}
      />

      <EmbeddedTable
        tagName={"Message"}
        tableType={"custom-vertical"}
        tableItems={[]}
      >
        <div
          class="whitespace-pre rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
        >
          {eventData?.message}
        </div>
      </EmbeddedTable>
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(eventData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
