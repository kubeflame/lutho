<script lang="ts">
  import CodeMirror from "../../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml } from "js-yaml";
  import {
    EmptyGVRK,
    EventV1GVRK,
    HelmChartV1GVRK,
    parseFieldSelector,
    randomUUID,
    routeString,
    tabs,
  } from "../../lib/util";
  import {
    type Helm,
    type HelmRepoData,
    type TabQueryParam,
  } from "../../lib/types";
  import HeaderElement from "../../lib/HeaderElement.svelte";
  import Tabs from "../../lib/Tabs.svelte";
  import ResourceToolbar from "../../lib/ResourceToolbar.svelte";
  import Details from "../../lib/Details.svelte";
  import socketStore from "../../lib/socketStore";
  import RouterPage from "../../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../../lib/ResourceToolbarBreadcrumbs.svelte";
  import EmbeddedTable from "../../lib/tables/EmbeddedTable.svelte";
  import SvgIcon from "../../lib/SvgIcon.svelte";
  import Events from "../../lib/Events.svelte";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml, tabs.events];
  const {
    sockError,
    sockState,
    isLoading,
    dataSend,
    dataGet,
    dataList,
    dataUpdate,
  } = socketStore();

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let helmReleaseData: Helm.Release;
  let codeMirrorChanged: boolean;
  let showSecondView: boolean;
  let helmRepos: HelmRepoData[] = JSON.parse(
    localStorage.getItem("helmRepos") as string,
  );
  let chartRefDialog: HTMLDialogElement;
  let chartVersion: string;
  let chartURL: string;
  let chartName: string = "";
  let chartIsOCI: boolean;
  let eventListData: any;

  $: toolbarContent = [
    { index: 0, name: "Helm List", url: routeString.helm + "?tab=List" },
    { index: 1, name: params.name },
  ];
  $: chartName = helmReleaseData?.chart?.metadata?.name as string;
  $: chartURL = getRepo(helmRepos, chartName)?.url;
  $: chartVersion = getRepo(helmRepos, chartName)?.latestVersion as string;
  $: chartIsOCI = getRepo(helmRepos, chartName)?.isOCI;
  $: sendGet = {
    opID: randomUUID(),
    type: "helmGet",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: EmptyGVRK,
    },
  } as any;
  $: $dataSend = [sendGet];
  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      name: params.name,
      kubeGVRK: EventV1GVRK,
      kubeOptions: {
        fieldSelector: parseFieldSelector({
          kind: HelmChartV1GVRK.kind,
          name: params.name,
          namespace: params.namespace,
        }),
      },
    },
  } as any;
  $: if (tabQueryParam === "events") {
    $dataSend = [sendList];
  }

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      helmReleaseData = dg.data;
    }
  });

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      eventListData = dl.data;
    }
  });

  dataSend.subscribe((ds) => {
    if (ds && $sockState.state === WebSocket.CLOSED) $sockState.refresh = true;
  });

  function getRepo(repoData: HelmRepoData[], name: string): HelmRepoData {
    if (repoData) {
      const repoIndex = repoData?.findIndex((r) => {
        return r.name === name;
      });
      return repoData[repoIndex];
    } else {
      return {} as HelmRepoData;
    }
  }

  function upgradeRelease() {
    if ($sockState.state === WebSocket.CLOSED) $sockState.refresh = true;
    $dataSend = [
      {
        type: "helmUpgrade",
        request: {
          namespace: params.namespace,
          name: params.name,
          helmOptions: {
            chartName: chartName,
            chartVersion: chartVersion,
            isOCI: chartIsOCI,
            repoURL: chartURL,
          },
          data: $docStore,
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
    onClickSubmit={() => chartRefDialog.showModal()}
  >
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
    <!-- <button
      slot="helmDefaultValues"
      class="btn btn-xs bg-base-100 outline-base-200 hover:btn-primary text-sm outline outline-1 drop-shadow-md"
      on:click={() => (showSecondView = !showSecondView)}
    >
      Default Values
    </button> -->
  </ResourceToolbar>
  {#if tabQueryParam === "details"}
    <Details title={"Release Information"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          { name: "Name", value: helmReleaseData?.name },
          { name: "Namespace", value: helmReleaseData?.namespace },
          {
            name: "Status",
            value: helmReleaseData?.info?.status,
          },
          {
            name: "Description",
            value: helmReleaseData?.info?.description,
          },
          {
            name: "First Deployed",
            value: helmReleaseData?.info?.first_deployed,
          },
          {
            name: "Last Deployed",
            value: helmReleaseData?.info?.last_deployed,
          },
          {
            name: "Notes",
            value: helmReleaseData?.info?.notes,
            className: "whitespace-pre-wrap",
          },
        ]}
      />
    </Details>

    <Details title={"Chart Information"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          { name: "Name", value: helmReleaseData?.chart?.metadata?.name },
          {
            name: "Version",
            value: helmReleaseData?.chart?.metadata?.version,
          },
          {
            name: "App Version",
            value: helmReleaseData?.chart?.metadata?.appVersion,
          },
          {
            name: "Description",
            value: helmReleaseData?.chart?.metadata?.description,
            className: "whitespace-pre-wrap",
          },
          {
            name: "Home",
            value: helmReleaseData?.chart?.metadata?.home,
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(helmReleaseData?.config)}
      bind:codeMirrorChanged
      bind:docStore
    />
    <dialog bind:this={chartRefDialog} class="modal modal-middle">
      <div
        class="modal-box border-base-100 bg-base-200 outline-base-200 grid max-h-[calc(70vh)]
          max-w-[calc(50vw)] grid-cols-1 rounded-lg border p-2 outline outline-1 drop-shadow-lg"
      >
        <div class="mb-5 flex items-center gap-x-1 text-sm">
          <SvgIcon classNames={"h-5 w-5"} strokeWidth={1.5} type={"update"} />
          Upgrade Release
        </div>

        <div class="grid grid-cols-1 gap-y-2">
          <div class="join h-6 w-full items-center">
            <label
              for="chart-name-input"
              class="join-item bg-base-300 h-full w-fit p-0.5 pl-2 pr-2 text-sm font-normal"
            >
              Chart Name
            </label>
            <input
              id="chart-name-input"
              bind:value={chartName}
              type="text"
              class="input input-xs join-item input-bordered bg-base-100 text-base-300
                flex grow focus:outline-0"
            />
          </div>
          <div class="join h-6 w-full items-center">
            <label
              for="chart-url-input"
              class="join-item bg-base-300 h-full w-fit p-0.5 pl-2 pr-2 text-sm font-normal"
            >
              Chart URL
            </label>
            <input
              id="chart-url-input"
              bind:value={chartURL}
              type="text"
              class="input input-xs join-item input-bordered bg-base-100 text-base-300
                flex grow focus:outline-0"
            />
          </div>
          <div class="join h-6 w-full items-center">
            <label
              for="chart-version-input"
              class="join-item bg-base-300 h-full w-fit p-0.5 pl-2 pr-2 text-sm font-normal"
            >
              Chart Version
            </label>
            <input
              id="chart-version-input"
              bind:value={chartVersion}
              type="text"
              class="input input-xs join-item input-bordered bg-base-100 flex grow focus:outline-0"
            />
          </div>
        </div>

        <div class="modal-action">
          <button
            class="btn btn-xs outline-success hover:bg-success focus:outline-success place-self-end
              outline outline-1 drop-shadow-md focus:outline-1"
            on:click={() => {
              upgradeRelease();
              chartRefDialog.close();
            }}
          >
            Upgrade
          </button>
          <button
            class="btn btn-xs outline-error hover:bg-error focus:outline-error place-self-end
              outline outline-1 drop-shadow-md focus:outline-1"
            on:click={() => chartRefDialog.close()}
          >
            Close
          </button>
        </div>
      </div>
    </dialog>
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
