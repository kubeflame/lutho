<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { onMount } from "svelte";
  import { routeString, DeploymentV1GVRK } from "../lib/util";
  import type { V1DeploymentList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  let deployListData: V1DeploymentList;

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "Deployment List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...DeploymentV1GVRK,
      },
    },
  ];

  $: deployListData = $dataList;
</script>

<HeaderElement>
  <NamespaceSelect slot="namespace" />
</HeaderElement>

<RouterPage>
  <ResourceToolbar slot="resource-toolbar" bind:toolbarContent />
  {#if $sockError}
    <ErrorPage bind:errorMessage={$sockError} />
  {:else if $isLoading}
    <LoadingNewton />
  {:else if deployListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each deployListData.items as deploy}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.deployList}/{deploy.metadata
                  ?.namespace}/{deploy.metadata?.name}"
                use:link
              >
                {deploy.metadata?.name}
              </a>
            </td>
            <td>{deploy.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</RouterPage>
