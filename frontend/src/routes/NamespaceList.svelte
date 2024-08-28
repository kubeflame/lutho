<script lang="ts">
  import { routeString, NamespaceV1GVRK, randomUUID } from "../lib/util";
  import type { V1NamespaceList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

  let namespaceListData: V1NamespaceList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Namespace List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      kubeGVRK: NamespaceV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      namespaceListData = dl.data;
    }
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          name: item.metadata?.name,
          kubeGVRK: NamespaceV1GVRK,
        },
      },
    ];

    dataDelete.subscribe((err) => {
      if (!err) $dataSend = [sendList];
    });
  }
</script>

<HeaderElement />

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.namespaceList}
    isNamespaced={false}
    tableHead={["Name", "Created At", ""]}
    items={namespaceListData?.items}
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
