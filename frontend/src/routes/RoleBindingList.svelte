<script lang="ts">
  import { namespace } from "../lib/stores";
  import { randomUUID, RoleBindingV1GVRK, routeString } from "../lib/util";
  import type { V1RoleBindingList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";

  let roleBindingListData: V1RoleBindingList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "RoleBinding List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      namespace: $namespace,
      kubeGVRK: RoleBindingV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      roleBindingListData = dl.data;
    }
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          name: item.metadata?.name,
          namespace: item.metadata?.namespace,
          kubeGVRK: RoleBindingV1GVRK,
        },
      },
    ];

    dataDelete.subscribe((err) => {
      if (!err) $dataSend = [sendList];
    });
  }
</script>

<HeaderElement showNamespaceSelect />

<RouterPage bind:errorMessage={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.roleBindingList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Created At", ""]}
    items={roleBindingListData?.items}
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
