<script lang="ts">
  import { link } from "svelte-spa-router";
  import ErrorPage from "../lib/ErrorPage.svelte";
  import LoadingNewton from "../lib/LoadingNewton.svelte";
  import { namespace } from "../lib/stores";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { RoleV1GVRK, routeString } from "../lib/util";
  import type { V1RoleList } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import { KubeDataOpType } from "../lib/types";
  import socketStore from "../lib/socketStore";

  let roleListData: V1RoleList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  onMount(async () => {
    $isLoading = true;
  });

  $: toolbarContent = [{ index: 0, name: "Role List" }];

  $: $dataSend = [
    {
      type: KubeDataOpType.list,
      request: {
        namespace: $namespace,
        ...RoleV1GVRK,
      },
    },
  ];

  $: roleListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: KubeDataOpType.list,
          request: {
            namespace: $namespace,
            ...RoleV1GVRK,
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
  {:else if roleListData}
    <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
        </tr>
      </thead>
      <tbody>
        {#each roleListData.items as role}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.roleList}/{role.metadata?.namespace}/{role
                  .metadata?.name}"
                use:link
              >
                {role.metadata?.name}
              </a>
            </td>
            <td>{role.metadata?.namespace}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>
