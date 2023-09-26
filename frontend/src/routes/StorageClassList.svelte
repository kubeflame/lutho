<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { StorageClassV1GVRK, routeString } from "../lib/util";
  import type { V1StorageClassList } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let storageClassListData: V1StorageClassList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "StorageClass List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...StorageClassV1GVRK,
      },
    },
  ];

  $: storageClassListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            ...StorageClassV1GVRK,
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
  {:else if storageClassListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        {#each storageClassListData.items as storageClass}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.storageClassList}/{storageClass.metadata
                  ?.name}"
                use:link
              >
                {storageClass.metadata?.name}
              </a>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
