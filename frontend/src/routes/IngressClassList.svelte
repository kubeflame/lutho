<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { routeString, IngressClassV1GVRK } from "../lib/util";
  import type { V1IngressClassList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let ingressClassListData: V1IngressClassList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "IngressClass List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...IngressClassV1GVRK,
      },
    },
  ];

  $: ingressClassListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            ...IngressClassV1GVRK,
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
  {:else if ingressClassListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
        </tr>
      </thead>
      <tbody>
        {#each ingressClassListData.items as ingressClass}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.ingressClassList}/{ingressClass.metadata
                  ?.name}"
                use:link
              >
                {ingressClass.metadata?.name}
              </a>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
