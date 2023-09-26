<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { routeString, NetworkPolicyV1GVRK } from "../lib/util";
  import type { V1NetworkPolicyList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let networkPolicyListData: V1NetworkPolicyList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "NetworkPolicy List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...NetworkPolicyV1GVRK,
      },
    },
  ];

  $: networkPolicyListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...NetworkPolicyV1GVRK,
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
  {:else if networkPolicyListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each networkPolicyListData.items as networkPolicy}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.networkPolicyList}/{networkPolicy.metadata
                  ?.namespace}/{networkPolicy.metadata?.name}"
                use:link
              >
                {networkPolicy.metadata?.name}
              </a>
            </td>
            <td>{networkPolicy.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
