<script lang="ts">
  import { StorageClassV1GVRK, randomUUID, routeString } from "../lib/util";
  import type { V1StorageClassList } from "@kubernetes/client-node/dist/gen/api";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";

  let storageClassListData: V1StorageClassList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "StorageClass List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      kubeGVRK: StorageClassV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      storageClassListData = dl.data;
    }
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

    dataDelete.subscribe((err) => {
      if (!err) $dataSend = [sendList];
    });
  }
</script>

<HeaderElement />

<RouterPage bind:errorMessage={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.storageClassList}
    isNamespaced={false}
    tableHead={["Name", "Created At", ""]}
    items={storageClassListData?.items}
    embeddedOptionsDataFn={(item) => [
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
</RouterPage>
