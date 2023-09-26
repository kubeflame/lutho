<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { routeString, ServiceV1GVRK } from "../lib/util";
  import type { V1ServiceList } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let serviceListData: V1ServiceList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "Service List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...ServiceV1GVRK,
      },
    },
  ];

  $: serviceListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...ServiceV1GVRK,
          },
        },
      ];
  });
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
  {:else if serviceListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each serviceListData.items as service}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.serviceList}/{service.metadata
                  ?.namespace}/{service.metadata?.name}"
                use:link
              >
                {service.metadata?.name}
              </a>
            </td>
            <td>{service.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
