<script lang="ts">
  import { namespace } from "../lib/stores";
  import { randomUUID, routeString, ServiceV1GVRK } from "../lib/util";
  import type { V1ServiceList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

  let serviceListData: V1ServiceList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Service List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      namespace: $namespace,
      kubeGVRK: ServiceV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      serviceListData = dl.data;
    }
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          name: item.metadata?.name,
          namespace: item.metadata?.namespace,
          kubeGVRK: ServiceV1GVRK,
        },
      },
    ];

    dataDelete.subscribe((err) => {
      if (!err) $dataSend = [sendList];
    });
  }
</script>

<HeaderElement showNamespaceSelect />

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.serviceList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Created At", ""]}
    items={serviceListData?.items}
  >
    <td
      class="flex place-items-center items-center justify-end"
      slot="embeddedOptions"
      let:item
    >
      <EmbeddedOptions
        embeddedOptionsData={[
          {
            fn: () => {},
            dialog: {
              action: () => onDelete(item),
              type: "Delete",
              resourceName: item.metadata?.name,
            },
            classes: "hover:btn-error",
            icon: "trash",
          },
        ]}
      />
    </td>
  </ListTable>
</RouterPage>
