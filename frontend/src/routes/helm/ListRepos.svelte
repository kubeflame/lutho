<script lang="ts">
  import ResourceToolbar from "../../lib/ResourceToolbar.svelte";
  import ResourceToolbarBreadcrumbs from "../../lib/ResourceToolbarBreadcrumbs.svelte";
  import RouterPage from "../../lib/RouterPage.svelte";
  import EmbeddedOptions from "../../lib/tables/EmbeddedOptions.svelte";
  import { EmptyGVRK, randomUUID, routeString } from "../../lib/util";
  import type {
    Alert,
    HelmChartTags,
    HelmRepoData,
    TabQueryParam,
  } from "../../lib/types";
  import { link } from "svelte-spa-router";
  import SvgIcon from "../../lib/SvgIcon.svelte";
  import socketStore from "../../lib/socketStore";
  import BottomAlert from "../../lib/BottomAlert.svelte";
  import { updatingRepos } from "../../lib/stores";

  export let tabQueryParam: TabQueryParam;

  const errorClasses: string = "input-error focus:border-error";
  const normalClasses: string = "focus:border-primary";
  const svgIconStrokeWidth: number = 1.5;
  const svgIconClassNames: string = "h-5 w-5";

  const { sockError, isLoading, dataSend, dataGet } = socketStore();

  let errorMessage: string = "";
  let loading: boolean = false;
  let newRepoDialog: HTMLDialogElement;
  let repoName: string = "";
  let repoURL: string = "";
  let isErrorRepoName: boolean = false;
  let isErrorRepoURL: boolean = false;
  let isErrorRepoType: boolean = false;
  let isOCI: boolean = false;
  let helmRepos: HelmRepoData[];

  $: alert = {
    message: null,
    type: "error",
    className: "flex mt-4 h-8",
  } as Alert;

  $: toolbarContent = [{ index: 0, name: "Helm Repository" }];
  $: helmRepos = JSON.parse(localStorage.getItem("helmRepos") as string) ?? [];

  $: sendGet = {
    opID: randomUUID(),
    type: "helmGetTags",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID && helmRepos) {
      const repoIndex = helmRepos.findIndex((r) => {
        return r.name === dg.data.chartName;
      });

      helmRepos[repoIndex].latestVersion = dg.data.chartTags[0];
      helmRepos[repoIndex].allVersions = dg.data.chartTags;
      localStorage.setItem("helmRepos", JSON.stringify(helmRepos));
      $updatingRepos[helmRepos[repoIndex].name] = false;
    }
  });

  async function newRepo(name: string) {
    let findRepoName: string = "";

    if (helmRepos) {
      findRepoName = helmRepos.find((repo: HelmRepoData) => {
        return repo.name === name;
      })?.name as string;
    }

    if (findRepoName === name) {
      alert.message = `Repository name '${repoName}' already exists, please specify a different name.`;
      isErrorRepoName = true;
    } else if (repoName === "" && repoURL === "") {
      alert.message = `Repository Name and URL must not be empty.`;
      isErrorRepoName = true;
      isErrorRepoURL = true;
    } else if (repoName === "") {
      alert.message = `Repository Name must not be empty.`;
      isErrorRepoName = true;
    } else if (repoURL === "") {
      alert.message = `Repository URL must not be empty.`;
      isErrorRepoURL = true;
    } else if (repoURL.startsWith("oci://") && !isOCI) {
      alert.message = `Repository Type was set to Helm, but the URL is OCI formatted.`;
      isErrorRepoType = true;
      isErrorRepoURL = true;
    } else if (
      (repoURL.startsWith("http://") || repoURL.startsWith("https://")) &&
      isOCI
    ) {
      alert.message = `Repository Type was set to OCI, but the URL is not OCI formatted.`;
      isErrorRepoType = true;
      isErrorRepoURL = true;
    } else {
      $dataSend = [
        {
          ...sendGet,
          request: {
            helmOptions: {
              chartName: name,
              repoURL: repoURL,
              isOCI: isOCI,
            },
          },
        },
      ];

      $updatingRepos[name] = true;

      helmRepos = [
        ...helmRepos,
        {
          name: repoName,
          url: repoURL,
          isOCI: isOCI,
        },
      ];

      localStorage.setItem("helmRepos", JSON.stringify(helmRepos));
      closeNewRepoDialog();
    }
  }

  function updateAll() {
    if (!helmRepos) return;

    helmRepos.forEach((r) => {
      $updatingRepos[r.name] = true;
      $dataSend = [
        {
          ...sendGet,
          request: {
            helmOptions: {
              chartName: r.name,
              repoURL: r.url,
              isOCI: r.isOCI,
            },
          },
        },
      ];
    });
  }

  function deleteRepo(repoName: string) {
    helmRepos = helmRepos.filter((repo: HelmRepoData) => {
      return repo.name !== repoName;
    });
    localStorage.setItem("helmRepos", JSON.stringify(helmRepos));
  }

  function closeNewRepoDialog() {
    isErrorRepoName = false;
    isErrorRepoURL = false;
    isErrorRepoType = false;
    repoName = "";
    repoURL = "";
    newRepoDialog.close();
  }
</script>

<RouterPage bind:errorMessage bind:loading>
  <ResourceToolbar slot="resource-toolbar" bind:tabQueryParam>
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
    <div class="mr-4 flex place-items-center space-x-2" slot="custom">
      <button
        class="btn btn-xs bg-base-100 text-sm outline outline-1 outline-primary drop-shadow-md hover:btn-primary"
        on:click={() => newRepoDialog.showModal()}
      >
        New
      </button>
      <button
        class="btn btn-xs bg-base-100 text-sm outline outline-1 outline-info drop-shadow-md hover:btn-info"
        on:click={updateAll}
      >
        Update All
      </button>
    </div>
  </ResourceToolbar>

  <dialog bind:this={newRepoDialog} class="modal modal-middle">
    <div
      class="modal-box grid max-h-[calc(70vh)] max-w-[calc(50vw)] grid-cols-1 rounded-lg border border-base-100 bg-base-200 p-2 outline outline-1 outline-base-200 drop-shadow-lg"
    >
      <div class="mb-5 flex items-center gap-x-1 text-sm">
        <SvgIcon
          classNames={svgIconClassNames}
          strokeWidth={svgIconStrokeWidth}
          type={"plusCircle"}
        />
        Add a new repository
      </div>

      <div class="grid grid-cols-1 gap-y-2">
        <div class="join h-6 w-full items-center gap-x-0.5">
          <label
            for="registry-type"
            class="join-item flex h-full w-fit items-center bg-base-300 p-1 pl-2 pr-2 text-sm font-normal"
          >
            Registry Type
          </label>
          <span id="registry-type" class="join rounded-l-md">
            <input
              class="btn join-item input-bordered btn-xs checked:shadow"
              type="radio"
              name="registry-type"
              aria-label="Helm"
              checked={!isOCI}
              on:click={() => {
                isOCI = false;
                isErrorRepoType = false;
                isErrorRepoURL = false;
                alert.message = null;
              }}
            />
            <input
              class="btn join-item input-bordered btn-xs checked:shadow"
              type="radio"
              name="registry-type"
              aria-label="OCI-based"
              checked={isOCI}
              on:click={() => {
                isOCI = true;
                isErrorRepoType = false;
                isErrorRepoURL = false;
                alert.message = null;
              }}
            />
          </span>
        </div>

        <div class="join h-6 w-full items-center">
          <label
            for="repo-name-input"
            class="join-item h-full w-14 bg-base-300 p-0.5 pl-2 pr-2 text-sm font-normal"
          >
            Name
          </label>
          <input
            id="repo-name-input"
            bind:value={repoName}
            type="text"
            class="input input-xs join-item input-bordered flex grow bg-base-100 focus:outline-0
              {isErrorRepoName ? errorClasses : normalClasses}"
            on:input={() => {
              isErrorRepoName = false;
              alert.message = null;
            }}
          />
        </div>
        <div class="join h-6 w-full items-center">
          <label
            for="repo-url-input"
            class="join-item h-full w-14 bg-base-300 p-0.5 pl-2 pr-2 text-sm font-normal"
          >
            URL
          </label>
          <input
            id="repo-url-input"
            bind:value={repoURL}
            type="text"
            class="input input-xs join-item input-bordered flex grow bg-base-100 focus:outline-0
              {isErrorRepoURL ? errorClasses : normalClasses}"
            on:input={() => {
              isErrorRepoURL = false;
              alert.message = null;
            }}
          />
        </div>
      </div>
      <div class="modal-action">
        <button
          class="btn btn-xs place-self-end outline outline-1 outline-success !drop-shadow-md hover:bg-success focus:outline-1 focus:outline-success"
          on:click={() => newRepo(repoName)}
          disabled={alert.message ? true : false}
        >
          Add
        </button>
        <button
          class="btn btn-xs place-self-end outline outline-1 outline-error drop-shadow-md hover:bg-error focus:outline-1 focus:outline-error"
          on:click={() => closeNewRepoDialog()}
        >
          Close
        </button>
      </div>

      <BottomAlert bind:alert />
    </div>
  </dialog>

  <table class="table table-pin-rows">
    <thead>
      <tr class="bg-base-200 shadow-sm">
        <th>Name</th>
        <th>URL</th>
        <th>Latest Version</th>
        <th />
      </tr>
    </thead>
    {#if helmRepos}
      <tbody>
        {#each helmRepos as repo}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.helm}/repository/{repo.name}"
                use:link
              >
                {repo.name}
              </a>
            </td>
            <td>{repo.url}</td>
            {#if $updatingRepos[repo.name]}
              <td>
                <span
                  class="loading loading-ring loading-sm flex place-self-center"
                />
              </td>
            {:else}
              <td>{repo.latestVersion}</td>
            {/if}
            <td class="flex place-items-center items-center justify-end">
              <EmbeddedOptions
                embeddedOptionsData={[
                  {
                    fn: () => {},
                    dialog: {
                      action: () => deleteRepo(repo.name),
                      type: "Delete",
                      resourceName: repo.name,
                    },
                    classes: "hover:btn-error",
                    icon: "trash",
                  },
                  {
                    fn: async () => {
                      $updatingRepos[repo.name] = true;

                      $dataSend = [
                        {
                          ...sendGet,
                          request: {
                            kubeGVRK: EmptyGVRK,
                            helmOptions: {
                              chartName: repo.name,
                              repoURL: repo.url,
                              isOCI: repo.isOCI,
                            },
                          },
                        },
                      ];
                    },
                    classes: "hover:btn-info",
                    icon: "update",
                  },
                ]}
              />
            </td>
          </tr>
        {/each}
      </tbody>
    {/if}
  </table>
</RouterPage>
