<script lang="ts">
  import { EmptyGVRK, randomUUID, routeString, tabs } from "../../lib/util";
  import HeaderElement from "../../lib/HeaderElement.svelte";
  import Tabs from "../../lib/Tabs.svelte";
  import RouterPage from "../../lib/RouterPage.svelte";
  import ResourceToolbar from "../../lib/ResourceToolbar.svelte";
  import ResourceToolbarBreadcrumbs from "../../lib/ResourceToolbarBreadcrumbs.svelte";
  import Details from "../../lib/Details.svelte";
  import EmbeddedTable from "../../lib/tables/EmbeddedTable.svelte";
  import socketStore from "../../lib/socketStore";
  import type { TabQueryParam } from "../../lib/types";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let tabQueryParam: TabQueryParam;
  let errorMessage: string = "";
  let loading: boolean = false;
  let repoInput: string = "";
  let repoChartData: any;

  $: toolbarContent = [
    {
      index: 0,
      name: "Helm Repository",
      url: routeString.helm + "?tab=Repository",
    },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "helmPull",
    request: {
      kubeGVRK: EmptyGVRK,
      helmOptions: {
        chartName: "kubernetes-dashboard",
        chartVersion: "6.0.7",
        envPath: repoInput,
        repoURL: "https://kubernetes.github.io/dashboard",
      },
    },
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      repoChartData = dg.data;
    }
  });

  function btnSend() {
    $dataSend = [sendGet];
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

<RouterPage bind:errorMessage bind:loading>
  <ResourceToolbar slot="resource-toolbar" bind:tabQueryParam>
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  {#if tabQueryParam === "details"}
    <input
      class="input input-xs input-bordered bg-base-100"
      bind:value={repoInput}
      type="text"
    />
    <button class="btn" on:click={btnSend}>REPO</button>

    Name: <br />
    Created: <br />
    Description: <br />
    Version: <br />
    AppVersion: <br />
    KubeVersion: <br />
    Dependencies: <br />
    Maintainers: <br />
    Sources: <br />
    URLs: <br />

    <Details title={"Resource Information"}>
      <EmbeddedTable
        tagName={"Summary"}
        tableType={"custom-vertical"}
        tableItems={[
          { name: "Name", value: "Placeholder Name" },
          { name: "Description", value: "Placeholder Description" },
          {
            name: "Creation Timestamp",
            value: "Placeholder Created",
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={[]}
      />

      <EmbeddedTable tableType={"badges"} tagName={"Labels"} tableItems={[]} />
    </Details>
  {/if}
</RouterPage>
