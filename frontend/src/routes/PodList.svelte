<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import {
    namespace,
    availableContainers,
    embeddedActiveTab,
  } from "../lib/stores";
  import { routeString, PodV1GVRK, icons } from "../lib/util";
  import EmbeddedOptions from "../lib/EmbeddedOptions.svelte";
  import type { V1Container, V1PodList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import { fade } from "svelte/transition";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType, TabIndex } from "../lib/types";
  import socketStore from "../lib/socketStore";
  import { onMount } from "svelte";

  let podListData: V1PodList;

  const getContainers = (containers: V1Container[]): string[] => {
    return containers.map((c) => {
      return c.name;
    });
  };

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "Pod List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...PodV1GVRK,
      },
    },
  ];

  $: podListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...PodV1GVRK,
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
  {:else if podListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
          <th>Status</th>
          <th />
        </tr>
      </thead>
      <tbody>
        {#each podListData.items as pod}
          <tr class="">
            <td class="flex items-center gap-x-2">
              <div
                class="badge {pod.status?.phase === 'Running'
                  ? 'badge-primary'
                  : 'badge-warning'} badge-xs"
              />
              <a
                class="hover:text-primary"
                href="{routeString.podList}/{pod.metadata?.namespace}/{pod
                  .metadata?.name}"
                use:link
              >
                {pod.metadata?.name}
              </a>
            </td>
            <td>{pod.metadata?.namespace}</td>
            <td>{pod.status?.phase}</td>
            <td class="flex place-items-center items-center">
              <EmbeddedOptions
                embeddedOptionsData={[
                  {
                    fn: () => {
                      $dataSend = [
                        {
                          type: KubeDataOpType.delete,
                          request: {
                            namespace: pod.metadata?.namespace,
                            name: pod.metadata?.name,
                            ...PodV1GVRK,
                          },
                        },
                      ];
                    },
                    classes: "hover:btn-error",
                    icon: icons.trash,
                  },
                  {
                    fn: () => {
                      $embeddedActiveTab = TabIndex.SHELL_EMBED;
                      $availableContainers = getContainers(
                        pod.spec?.containers ?? [],
                      );
                    },
                    classes: "hover:btn-warning",
                    icon: icons.cmdLine,
                    url: `${routeString.podList}/${pod.metadata?.namespace}/${pod.metadata?.name}`,
                    hide: pod.status?.phase !== "Running",
                  },
                ]}
              />
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
