<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { fade } from "svelte/transition";
  import { routeString, NodeV1GVRK } from "../lib/util";
  import type { V1NodeCondition, V1NodeList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let nodeListData: V1NodeList;

  const getNodeStatus = (nc: V1NodeCondition[]): string => {
    return (
      nc.find(
        (item: V1NodeCondition): boolean =>
          item.type === "Ready" && item.status === "True",
      )?.type || "Not Ready"
    );
  };

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Node List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        ...NodeV1GVRK,
      },
    },
  ];

  $: nodeListData = $dataList;
</script>

<HeaderElement>
  <ResourceToolbar slot="toolbar" bind:toolbarContent />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if nodeListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {#each nodeListData.items as node}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.nodeList}/{node.metadata?.name}"
                use:link
              >
                {node.metadata?.name}
              </a>
              {#if node.spec?.unschedulable}
                <span
                  class="badge badge-ghost badge-outline badge-sm font-light tracking-wide"
                >
                  SchedulingDisabled
                </span>
              {/if}
            </td>
            <td>
              {getNodeStatus(node.status?.conditions ?? [])}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
