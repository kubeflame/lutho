<script lang="ts">
  import { link } from "svelte-spa-router";
  import { namespace } from "../../lib/stores";
  import { EmptyGVRK, randomUUID, routeString } from "../../lib/util";
  import { type Helm } from "../../lib/types";
  import ResourceToolbar from "../../lib/ResourceToolbar.svelte";
  import socketStore from "../../lib/socketStore";
  import RouterPage from "../../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../../lib/ResourceToolbarBreadcrumbs.svelte";
  import EmbeddedOptions from "../../lib/tables/EmbeddedOptions.svelte";

  let helmReleaseListData: Helm.Release[];
  let dryRun: boolean = false;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Helm List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "helmList",
    request: {
      namespace: $namespace,
      kubeGVRK: EmptyGVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      helmReleaseListData = dl.data;
    }
  });

  dataDelete.subscribe((err) => {
    if (!err) $dataSend = [sendList];
  });

  function onDelete(release: Helm.Release) {
    $dataSend = [
      {
        type: "helmUninstall",
        request: {
          namespace: release.namespace,
          name: release.name,
          kubeGVRK: EmptyGVRK,
          helmOptions: {
            dryRun,
          },
        },
      },
    ];
  }
</script>

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>
  {#if helmReleaseListData}
    <table class="table-pin-rows table">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
          <th>Status</th>
          <th>Chart Version</th>
          <th>App Version</th>
          <th>Last Deployed</th>
          <th />
        </tr>
      </thead>
      <tbody>
        {#each helmReleaseListData as release}
          <tr>
            <td>
              <a
                class="hover:text-primary"
                href="{routeString.helm}/releases/{release?.namespace}/{release?.name}"
                use:link
              >
                {release?.name}
              </a>
            </td>
            <td>{release?.namespace}</td>
            <td>{release?.info?.status}</td>
            <td>{release?.chart?.metadata?.version}</td>
            <td>{release?.chart?.metadata?.appVersion}</td>
            <td>{release?.info?.last_deployed}</td>
            <td class="flex place-items-center items-center justify-end">
              <EmbeddedOptions
                embeddedOptionsData={[
                  {
                    fn: () => {},
                    dialog: {
                      action: () => onDelete(release),
                      type: "Delete",
                      resourceName: release?.name,
                    },
                    classes: "hover:btn-error",
                    icon: "trash",
                  },
                ]}
              />
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</RouterPage>
