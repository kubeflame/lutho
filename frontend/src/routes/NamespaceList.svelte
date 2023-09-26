<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { routeString, NamespaceV1GVRK } from "../lib/util";
  import type { V1NamespaceList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let namespaceListData: V1NamespaceList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "Namespace List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...NamespaceV1GVRK,
      },
    },
  ];

  $: namespaceListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            ...NamespaceV1GVRK,
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
  {:else if namespaceListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        {#each namespaceListData.items as namespace}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.namespaceList}/{namespace.metadata?.name}"
                use:link
              >
                {namespace.metadata?.name}
              </a>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
