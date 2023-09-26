<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { routeString, StatefulSetV1GVRK } from "../lib/util";
  import type { V1StatefulSetList } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let statefulSetListData: V1StatefulSetList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "StatefulSet List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...StatefulSetV1GVRK,
      },
    },
  ];

  $: statefulSetListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...StatefulSetV1GVRK,
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
  {:else if statefulSetListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each statefulSetListData.items as satefulSet}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.statefulSetList}/{satefulSet.metadata
                  ?.namespace}/{satefulSet.metadata?.name}"
                use:link
              >
                {satefulSet.metadata?.name}
              </a>
            </td>
            <td>{satefulSet.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
