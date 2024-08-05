<script lang="ts">
  import { scale } from "svelte/transition";
  import CodeMirror from "../../lib/codemirror/CodeMirror.svelte";
  import ResourceToolbar from "../../lib/ResourceToolbar.svelte";
  import RouterPage from "../../lib/RouterPage.svelte";
  import socketStore from "../../lib/socketStore";
  import type {
    DialogData,
    HelmRepoData,
    TabQueryParam,
  } from "../../lib/types";
  import SvgIcon from "../../lib/SvgIcon.svelte";
  import DialogElement from "../../lib/DialogElement.svelte";

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

  $: allVersions = (chartRepo?.allVersions as string[]) || [
    chartRepo?.latestVersion,
  ];
  $: helmShowValues = $showGet;

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

  function onShowValues() {
    if ($showState.state === WebSocket.CLOSED) $showState.refresh = true;

    $showSend = [
      {
        type: "helmShowValues",
        request: {
          helmOptions: {
            chartName: chartRepo.name,
            chartVersion: chartVersion,
            isOCI: chartRepo.isOCI,
            repoURL: chartRepo.url,
          },
        },
      },
    ];
  }

  function onInstall() {
    if ($installState.state === WebSocket.CLOSED) $installState.refresh = true;

    $installSend = [
      {
        type: "helmInstall",
        request: {
          namespace: releaseNamespace,
          name: releaseName,
          helmOptions: {
            chartName: chartRepo.name,
            chartVersion: chartVersion,
            isOCI: chartRepo.isOCI,
            repoURL: chartRepo.url,
          },
          data: $docStore,
        },
      },
    ];
  }

  $: dialogData = {
    action: onInstall,
    resourceName: releaseName,
    type: "Install",
    icon: "questionMark",
  } as DialogData;
</script>

<DialogElement bind:dialogData bind:showDialog />

<RouterPage bind:error={$showError} bind:loading={$showLoading}>
  <ResourceToolbar slot="resource-toolbar" bind:tabQueryParam>
    <div
      slot="custom"
      class="flex h-8 w-full place-items-center items-center
        justify-between gap-x-2 rounded-b-2xl rounded-t-none border-0 focus:outline-0"
    >
      <div class="join ml-4 min-w-fit items-center text-sm" in:scale>
        <label
          for="repo-select"
          class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
        >
          Repo
        </label>
        <select
          bind:value={chartRepo}
          id="repo-select"
          class="join-item select select-bordered select-xs bg-base-100
            focus:border-primary/60 rounded-lg pl-4 text-sm font-normal focus:outline-0"
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
          class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
        >
          Chart Version
        </label>
        <select
          bind:value={chartVersion}
          id="chart-version-select"
          class="join-item select select-bordered select-xs bg-base-100
            focus:border-primary/60 rounded-lg pl-4 text-sm font-normal focus:outline-0"
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
        class="btn btn-xs tooltip tooltip-bottom tooltip-primary bg-base-100 hover:bg-base-200
          join-item border-base-300 place-items-center items-center rounded-lg font-medium"
        data-tip="get the default chart values"
      >
        <SvgIcon classNames={"size-4"} type={"code"} />
      </button>

      <div class="join w-full items-center text-sm" in:scale>
        <label
          for="name-input"
          class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
        >
          Name
        </label>
        <input
          id="name-input"
          bind:value={releaseName}
          type="text"
          class="input input-xs join-item input-bordered bg-base-100
            focus:border-primary/60 flex grow text-sm focus:outline-0"
        />
      </div>

      <div class="join w-full items-center text-sm" in:scale>
        <label
          for="namespace-input"
          class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
        >
          Namespace
        </label>
        <input
          id="namespace-input"
          bind:value={releaseNamespace}
          type="text"
          class="input input-xs join-item input-bordered bg-base-100
            focus:border-primary/60 flex grow text-sm focus:outline-0"
        />
      </div>

      <div class="join">
        <button
          disabled={releaseName && releaseNamespace && chartRepo && chartVersion
            ? false
            : true}
          class="btn btn-sm bg-base-100 hover:bg-base-200 join-item
            place-items-center items-center rounded-br-2xl rounded-tr-none pr-5"
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
