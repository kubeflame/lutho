<script lang="ts">
  import { namespace } from "../lib/stores";
  import { routeString, ConfigMapV1GVRK } from "../lib/util";
  import type { V1ConfigMapList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

  let configMapListData: V1ConfigMapList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "ConfigMap List" }];

  $: $dataSend = [
    {
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: ConfigMapV1GVRK,
      },
    },
  ];

  $: configMapListData = $dataList;

  dataDelete.subscribe((err) => {
    if (!err)
      $dataSend = [
        {
          type: "list",
          request: {
            namespace: $namespace,
            kubeGVRK: ConfigMapV1GVRK,
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
          kubeGVRK: ConfigMapV1GVRK,
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
    hrefRoot={routeString.configMapList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Created At", ""]}
    items={configMapListData?.items}
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
