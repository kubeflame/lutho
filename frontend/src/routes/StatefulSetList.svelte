<script lang="ts">
  import { namespace } from "../lib/stores";
  import { randomUUID, routeString, StatefulSetV1GVRK } from "../lib/util";
  import type { V1StatefulSetList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

  let statefulSetListData: V1StatefulSetList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "StatefulSet List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      namespace: $namespace,
      kubeGVRK: StatefulSetV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      statefulSetListData = dl.data;
    }
  });

  dataDelete.subscribe((err) => {
    if (!err) $dataSend = [sendList];
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          name: item.metadata?.name,
          namespace: item.metadata?.namespace,
          kubeGVRK: StatefulSetV1GVRK,
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
    hrefRoot={routeString.statefulSetList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Created At", ""]}
    items={statefulSetListData?.items}
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
