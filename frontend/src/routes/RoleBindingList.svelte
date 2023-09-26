<script lang="ts">
  import { namespace } from "../lib/stores";
  import { RoleBindingV1GVRK, routeString } from "../lib/util";
  import type { V1RoleBindingList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

  let roleBindingListData: V1RoleBindingList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "RoleBinding List" }];

  $: $dataSend = [
    {
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: RoleBindingV1GVRK,
      },
    },
  ];

  $: roleBindingListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: "list",
          request: {
            namespace: $namespace,
            kubeGVRK: RoleBindingV1GVRK,
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
          namespace: item.metadata?.namespace,
          kubeGVRK: RoleBindingV1GVRK,
        },
      },
    ];
  }
</script>

<HeaderElement showNamespaceSelect />

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.roleBindingList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Created At", ""]}
    items={roleBindingListData?.items}
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
