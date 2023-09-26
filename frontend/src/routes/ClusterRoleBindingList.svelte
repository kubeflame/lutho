<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { ClusterRoleBindingV1GVRK, routeString } from "../lib/util";
  import type { V1ClusterRoleBindingList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let clusterRoleBindingListData: V1ClusterRoleBindingList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "ClusterRoleBinding List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...ClusterRoleBindingV1GVRK,
      },
    },
  ];

  $: clusterRoleBindingListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            ...ClusterRoleBindingV1GVRK,
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
  {:else if clusterRoleBindingListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        {#each clusterRoleBindingListData.items as clusterRoleBinding}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.clusterRoleBindingList}/{clusterRoleBinding
                  .metadata?.name}"
                use:link
              >
                {clusterRoleBinding.metadata?.name}
              </a>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
