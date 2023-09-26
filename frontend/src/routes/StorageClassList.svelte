<script lang="ts">
  import { StorageClassV1GVRK, routeString } from "../lib/util";
  import type { V1StorageClassList } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

  let storageClassListData: V1StorageClassList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "StorageClass List" }];

  $: $dataSend = [
    {
      type: "list",
      request: {
        kubeGVRK: StorageClassV1GVRK,
      },
    },
  ];

  $: storageClassListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: "list",
          request: {
            kubeGVRK: StorageClassV1GVRK,
          },
        },
      ];
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          name: item.metadata?.name,
          kubeGVRK: StorageClassV1GVRK,
        },
      },
    ];
  }
</script>

<HeaderElement />

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.storageClassList}
    isNamespaced={false}
    tableHead={["Name", "Created At", ""]}
    items={storageClassListData?.items}
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
