<script lang="ts">
  import { scale } from "svelte/transition";
  import CodeMirror from "../../lib/codemirror/CodeMirror.svelte";
  import ResourceToolbar from "../../lib/ResourceToolbar.svelte";
  import RouterPage from "../../lib/RouterPage.svelte";
  import socketStore from "../../lib/socketStore";
  import type {
    Alert,
    DialogData,
    DialogStep,
    HelmRepoData,
    TabQueryParam,
  } from "../../lib/types";
  import SvgIcon from "../../lib/SvgIcon.svelte";
  import DialogElement from "../../lib/DialogElement.svelte";
  import { randomUUID, routeString } from "../../lib/util";
  import { link } from "svelte-spa-router";
  import EmbeddedTable from "../../lib/tables/EmbeddedTable.svelte";
  import Details from "../../lib/Details.svelte";

  export let tabQueryParam: TabQueryParam;

  let docStore: any;
  let releaseName: string;
  let releaseNamespace: string;
  let chartRepo: HelmRepoData;
  let chartVersion: string;
  let helmRepoList =
    JSON.parse(localStorage.getItem("helmRepos") as string) ?? [];
  let allVersions: string[] = [];
  let showDialog: boolean = false;
  let helmShowValues: any;
  let alert: Alert = {
    message: null,
    type: null,
    className: "rounded-lg mt-4 p-1 pl-2",
  };
  let step: DialogStep;

  $: allVersions = (chartRepo?.allVersions as string[]) || [
    chartRepo?.latestVersion,
  ];

  $: sendInstall = {
    opID: randomUUID(),
    type: "helmInstall",
  } as any;

  $: sendShow = {
    opID: randomUUID(),
    type: "helmShowValues",
  } as any;

  const {
    sockError: installError,
    sockState: installState,
    isLoading: installLoading,
    dataSend: installSend,
    dataGet: installGet,
  } = socketStore();

  const {
    sockError: showError,
    sockState: showState,
    isLoading: showLoading,
    dataSend: showSend,
    dataGet: showGet,
  } = socketStore();

  showGet.subscribe((sg) => {
    if (sg && sg.op?.opID === sendShow.opID) {
      helmShowValues = sg.data;
    }
  });

  function onShowValues() {
    if ($showState.state === WebSocket.CLOSED) $showState.refresh = true;

    helmShowValues = "";

    $showSend = [
      {
        ...sendShow,
        request: {
          helmOptions: {
            chartName: chartRepo?.name,
            chartVersion: chartVersion,
            isOCI: chartRepo?.isOCI,
            repoURL: chartRepo?.url,
          },
        },
      },
    ];
  }

  function onInstall() {
    if ($installState.state === WebSocket.CLOSED) $installState.refresh = true;

    alert.message = null;
    alert.type = null;
    alert.timeout = null;
    $installError = "";
    $installGet = "";

    $installSend = [
      {
        ...sendInstall,
        request: {
          namespace: releaseNamespace,
          name: releaseName,
          helmOptions: {
            chartName: chartRepo?.name,
            chartVersion: chartVersion,
            isOCI: chartRepo?.isOCI,
            repoURL: chartRepo?.url,
          },
          data: $docStore,
        },
      },
    ];

    installError.subscribe((ie) => {
      if (ie) {
        alert.message = ie;
        alert.type = "error";
      }
    });

    installGet.subscribe((ig) => {
      if (ig) {
        step = "afterAction";
        alert.message = ig.data?.info.description;
        alert.type = "success";
        alert.timeout = 5;
      }
    });
  }

  $: dialogData = {
    action: onInstall,
    resourceName: releaseName,
    type: "Install",
    icon: "questionMark",
  } as DialogData;
</script>

<DialogElement
  bind:dialogData
  bind:showDialog
  bind:alert
  onClose={() => {
    alert.message = null;
    alert.type = null;
    alert.timeout = null;
    step = "action";
    $installGet = "";
  }}
  isLoading={$installLoading}
  bind:step
>
  <div slot="dialogContent" class="space-y-2 py-2">
    <EmbeddedTable
      tagName={"Summary"}
      tableType={"custom-vertical"}
      tableItems={[
        { name: "Status", value: $installGet?.data?.info?.status },
        { name: "Description", value: $installGet?.data?.info?.description },
      ]}
    />

    <Details title={"Release Notes"} collapseChecked={false}>
      <div class="overflow-hidden break-words rounded-lg font-mono text-sm">
        <div class="whitespace-pre-wrap">
          {$installGet?.data?.info?.notes}
        </div>
      </div>
    </Details>
  </div>

  <div slot="afterActionLinks" class="space-x-2">
    <a
      role="button"
      class="btn btn-xs outline outline-1 outline-base-200 drop-shadow-sm
        hover:bg-success focus:outline-1 focus:outline-success"
      href="{routeString.helm}/releases/{releaseNamespace}/{releaseName}"
      use:link
    >
      <SvgIcon type={"link"} />
      {releaseName}
    </a>
    <a
      role="button"
      class="btn btn-xs outline outline-1 outline-base-200 drop-shadow-sm
        hover:bg-success focus:outline-1 focus:outline-success"
      href="{routeString.helm}?tab=list"
      use:link
    >
      <SvgIcon type={"link"} />
      Helm List
    </a>
  </div>
</DialogElement>

<RouterPage bind:errorMessage={$showError} bind:loading={$showLoading}>
  <ResourceToolbar slot="resource-toolbar" bind:tabQueryParam>
    <div
      slot="custom"
      class="flex h-8 w-full place-items-center items-center
        justify-between gap-x-2 rounded-b-2xl rounded-t-none border-0 focus:outline-0"
    >
      <div class="join ml-4 min-w-fit items-center text-sm" in:scale>
        <label
          for="repo-select"
          class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
        >
          Repo
        </label>
        <select
          bind:value={chartRepo}
          id="repo-select"
          class="join-item select select-bordered select-xs rounded-lg
            bg-base-100 pl-4 text-sm font-normal focus:border-primary/60 focus:outline-0"
        >
          <option selected value="" />
          {#each helmRepoList as repo}
            <option value={repo}>{repo.name}</option>
          {/each}
        </select>
      </div>

      <div class="join min-w-fit items-center text-sm" in:scale>
        <label
          for="chart-version-select"
          class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
        >
          Chart Version
        </label>
        <select
          bind:value={chartVersion}
          id="chart-version-select"
          class="join-item select select-bordered select-xs rounded-lg
            bg-base-100 pl-4 text-sm font-normal focus:border-primary/60 focus:outline-0"
        >
          <option selected value="" />
          {#if chartRepo}
            {#each allVersions as version}
              <option value={version}>{version}</option>
            {/each}
          {/if}
        </select>
      </div>

      <button
        disabled={chartRepo && chartVersion ? false : true}
        on:click={onShowValues}
        class="btn join-item btn-xs tooltip tooltip-bottom tooltip-primary place-items-center
          items-center rounded-lg border-base-300 bg-base-100 font-medium hover:bg-base-200"
        data-tip="get the default chart values"
      >
        <SvgIcon classNames={"size-4"} type={"code"} />
      </button>

      <div class="join w-full items-center text-sm" in:scale>
        <label
          for="name-input"
          class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
        >
          Name
        </label>
        <input
          id="name-input"
          bind:value={releaseName}
          type="text"
          class="input input-xs join-item input-bordered flex
            grow bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>

      <div class="join w-full items-center text-sm" in:scale>
        <label
          for="namespace-input"
          class="join-item min-w-fit bg-base-200 p-0.5 pl-2 pr-2 font-normal"
        >
          Namespace
        </label>
        <input
          id="namespace-input"
          bind:value={releaseNamespace}
          type="text"
          class="input input-xs join-item input-bordered flex
            grow bg-base-100 text-sm focus:border-primary/60 focus:outline-0"
        />
      </div>

      <div class="join">
        <button
          disabled={releaseName && releaseNamespace && chartRepo && chartVersion
            ? false
            : true}
          class="btn join-item btn-sm place-items-center items-center
            rounded-br-2xl rounded-tr-none bg-base-100 pr-5 hover:bg-base-200"
          on:click={() => {
            showDialog = true;
          }}
        >
          <SvgIcon classNames={"size-5"} type={"arrowRight"} />
        </button>
      </div>
    </div>
  </ResourceToolbar>

  <CodeMirror value={helmShowValues} bind:docStore />
</RouterPage>
