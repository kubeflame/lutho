<script lang="ts">
  import { type KubeGVRK, type TabQueryParam } from "./types";
  import socketStore from "./socketStore";
  import { link } from "svelte-spa-router";
  import RouterPage from "./RouterPage.svelte";
  import ResourceToolbar from "./ResourceToolbar.svelte";
  import SvgIcon from "./SvgIcon.svelte";
  import { scale } from "svelte/transition";
  import {
    ClusterRoleBindingV1GVRK,
    ClusterRoleV1GVRK,
    ConfigMapV1GVRK,
    CronJobV1GVRK,
    CustomResourceDefinitionV1GVRK,
    DaemonSetV1GVRK,
    DeploymentV1GVRK,
    EventV1GVRK,
    IngressClassV1GVRK,
    IngressV1GVRK,
    JobV1GVRK,
    NamespaceV1GVRK,
    NetworkPolicyV1GVRK,
    NodeV1GVRK,
    PersistentVolumeClaimV1GVRK,
    PersistentVolumeV1GVRK,
    PodV1GVRK,
    ReplicaSetV1GVRK,
    ReplicationControllerV1GVRK,
    RoleBindingV1GVRK,
    RoleV1GVRK,
    SecretV1GVRK,
    ServiceAccountV1GVRK,
    ServiceV1GVRK,
    StatefulSetV1GVRK,
    StorageClassV1GVRK,
  } from "./util";

  export let tabQueryParam: TabQueryParam;

  const kubeResourceTypeList: KubeGVRK[] = [
    NamespaceV1GVRK,
    NodeV1GVRK,
    PersistentVolumeV1GVRK,
    PersistentVolumeClaimV1GVRK,
    ReplicationControllerV1GVRK,
    PodV1GVRK,
    ServiceV1GVRK,
    SecretV1GVRK,
    ServiceAccountV1GVRK,
    ConfigMapV1GVRK,
    DeploymentV1GVRK,
    DaemonSetV1GVRK,
    ReplicaSetV1GVRK,
    StatefulSetV1GVRK,
    EventV1GVRK,
    JobV1GVRK,
    CronJobV1GVRK,
    IngressV1GVRK,
    IngressClassV1GVRK,
    NetworkPolicyV1GVRK,
    RoleV1GVRK,
    RoleBindingV1GVRK,
    ClusterRoleV1GVRK,
    ClusterRoleBindingV1GVRK,
    StorageClassV1GVRK,
    CustomResourceDefinitionV1GVRK,
  ];
  const svgIconClassNames = "h-5 w-5";
  const { sockState, isLoading, dataSend, sockError, dataList } = socketStore();

  let resourceType: KubeGVRK;
  let searchDataResponse: any;
  let isNamespaced: boolean;
  let namespaceInput: string;
  let nameInput: string;
  let labelInput: string;
  let fieldInput: string;
  let showFilter: boolean;

  sockState.subscribe((s) => {
    if (s.bound) {
      isLoading.update((load) => (load = false));
    }
  });

  $: searchDataResponse = $dataList;
  $: if (!showFilter) (labelInput = ""), (fieldInput = "");

  function onSearch() {
    isNamespaced = false;
    $dataSend = [
      {
        type: "list",
        request: {
          namespace: resourceType?.isNamespaced ? namespaceInput : "",
          kubeGVRK: resourceType,
          kubeOptions: {
            fieldSelector: nameInput
              ? fieldInput.concat(",", `metadata.name=${nameInput}`)
              : fieldInput,
            labelSelector: labelInput,
          },
        },
      },
    ];
    if (resourceType.isNamespaced) isNamespaced = true;
  }
</script>

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar" bind:tabQueryParam>
    <div
      slot="custom"
      class="flex h-8 w-full place-items-center items-center justify-between gap-x-2
        rounded-b-2xl rounded-t-none border-0 focus:outline-0"
    >
      {#if tabQueryParam === "search"}
        <button
          on:click={() => (showFilter = !showFilter)}
          class="btn btn-sm tooltip tooltip-right tooltip-primary bg-base-100 hover:bg-base-200
          place-items-center items-center rounded-bl-2xl rounded-tl-none pl-5 font-medium"
          data-tip="more filters"
        >
          <SvgIcon classNames={svgIconClassNames} type={"filter"} />
        </button>

        <div class="join min-w-fit items-center text-sm" in:scale>
          <label
            for="resource-select"
            class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
          >
            Resource Type
          </label>
          <select
            bind:value={resourceType}
            id="resource-select"
            class="join-item select select-bordered select-xs bg-base-100 focus:border-primary/60
            rounded-lg pl-4 text-sm font-normal focus:outline-0"
          >
            <option selected disabled value="">Choose a resource</option>
            {#each kubeResourceTypeList as resource}
              <option value={resource}>{resource.kind}</option>
            {/each}
          </select>
        </div>

        <div class="join w-full items-center text-sm" in:scale>
          <label
            for="name-input"
            class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
          >
            Name
          </label>
          <input
            id="name-input"
            bind:value={nameInput}
            type="text"
            class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
            flex grow text-sm focus:outline-0"
          />
        </div>

        {#if resourceType?.isNamespaced}
          <div class="join w-full items-center text-sm" in:scale>
            <label
              for="namespace-input"
              class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
            >
              Namespace
            </label>
            <input
              id="namespace-input"
              bind:value={namespaceInput}
              type="text"
              class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
            flex grow text-sm focus:outline-0"
            />
          </div>
        {/if}

        {#if showFilter}
          <div class="join w-full items-center text-sm" in:scale>
            <label
              for="field-input"
              class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
            >
              Field Selector
            </label>
            <input
              id="field-input"
              bind:value={fieldInput}
              type="text"
              class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
              flex grow text-sm focus:outline-0"
            />
          </div>

          <div class="join h-6 w-full items-center text-sm" in:scale>
            <label
              for="label-input"
              class="join-item bg-base-200 min-w-fit p-0.5 pl-2 pr-2 font-normal"
            >
              Label Selector
            </label>
            <input
              id="label-input"
              bind:value={labelInput}
              type="text"
              class="input input-xs join-item input-bordered bg-base-100 focus:border-primary/60
              flex grow text-sm focus:outline-0"
            />
          </div>
        {/if}

        <button
          disabled={resourceType ? false : true}
          class="btn btn-sm bg-base-100 hover:bg-base-200 place-items-center
          items-center rounded-br-2xl rounded-tr-none pr-5"
          on:click={onSearch}
        >
          <SvgIcon classNames={svgIconClassNames} type={"arrowRight"} />
        </button>
      {/if}
    </div>
  </ResourceToolbar>

  {#if searchDataResponse}
    <table class="table-pin-rows table">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          {#if isNamespaced}
            <th>Namespace</th>
          {/if}
        </tr>
      </thead>
      <tbody>
        {#each searchDataResponse.items as searchItem}
          <tr class="">
            {#if isNamespaced}
              <td class="flex items-center gap-x-2">
                <a
                  class="hover:text-primary"
                  href="/{resourceType?.resource}/{searchItem.metadata
                    ?.namespace}/{searchItem.metadata?.name}"
                  use:link
                >
                  {searchItem.metadata?.name}
                </a>
              </td>
              <td>{searchItem.metadata?.namespace}</td>
            {:else}
              <td class="flex items-center gap-x-2">
                <a
                  class="hover:text-primary"
                  href="/{resourceType?.resource}/{searchItem.metadata?.name}"
                  use:link
                >
                  {searchItem.metadata?.name}
                </a>
              </td>
            {/if}
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</RouterPage>
