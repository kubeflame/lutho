<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { DaemonSetV1GVRK, routeString } from "../lib/util";
  import type { V1DaemonSetList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let daemonSetListData: V1DaemonSetList;

  onMount(async () => {
    $isLoading = true;
  });

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "DaemonSet List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...DaemonSetV1GVRK,
      },
    },
  ];

  $: daemonSetListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...DaemonSetV1GVRK,
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
  {:else if daemonSetListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each daemonSetListData.items as daemonSet}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.daemonSetList}/{daemonSet.metadata
                  ?.namespace}/{daemonSet.metadata?.name}"
                use:link
              >
                {daemonSet.metadata?.name}
              </a>
            </td>
            <td>{daemonSet.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
