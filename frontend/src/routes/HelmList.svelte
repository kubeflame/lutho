<script lang="ts">
  import { link } from "svelte-spa-router";
  import { namespace } from "../lib/stores";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import { fade } from "svelte/transition";
  import { routeString } from "../lib/util";
  import { KubeDataOpType, type Helm } from "../lib/types";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";

  let helmReleaseListData: Helm.Release[];

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Helm List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.helmList,
      request: {
        namespace: $namespace,
      },
    },
  ];

  $: helmReleaseListData = $dataList;
</script>

<HeaderElement>
  <NamespaceSelect slot="namespace" />
  <ResourceToolbar slot="toolbar" bind:toolbarContent />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if helmReleaseListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
          <th>Status</th>
          <th>Last Deployed</th>
        </tr>
      </thead>
      <tbody>
        {#each helmReleaseListData as release}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.helm}/{release.namespace}/{release.name}"
                use:link
              >
                {release.name}
              </a>
            </td>
            <td>{release.namespace}</td>
            <td>{release.info.status}</td>
            <td>{release.info.last_deployed}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
