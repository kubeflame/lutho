<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { routeString, ClusterRoleV1GVRK } from "../lib/util";
  import type { V1ClusterRoleList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let clusterRoleListData: V1ClusterRoleList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "ClusterRole List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...ClusterRoleV1GVRK,
      },
    },
  ];

  $: clusterRoleListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            ...ClusterRoleV1GVRK,
          },
        },
      ];
  });
</script>

<HeaderElement>
  <ResourceToolbar slot="toolbar" bind:toolbarContent />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if clusterRoleListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        {#each clusterRoleListData.items as clusterRole}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.clusterRoleList}/{clusterRole?.metadata
                  ?.name}"
                use:link
              >
                {clusterRole?.metadata?.name}
              </a>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
