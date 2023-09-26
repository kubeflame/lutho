<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { routeString, PersistentVolumeV1GVRK } from "../lib/util";
  import type { V1PersistentVolumeList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let persistentVolumeListData: V1PersistentVolumeList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "PersistentVolume List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...PersistentVolumeV1GVRK,
      },
    },
  ];

  $: persistentVolumeListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...PersistentVolumeV1GVRK,
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
  {:else if persistentVolumeListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each persistentVolumeListData.items as persistentVolume}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.persistentVolumeList}/{persistentVolume
                  .metadata?.namespace}/{persistentVolume.metadata?.name}"
                use:link
              >
                {persistentVolume.metadata?.name}
              </a>
            </td>
            <td>{persistentVolume.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
