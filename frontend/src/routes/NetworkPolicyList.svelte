<script lang="ts">
  import { namespace } from "../lib/stores";
  import { routeString, NetworkPolicyV1GVRK, randomUUID } from "../lib/util";
  import type { V1NetworkPolicyList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";

  let networkPolicyListData: V1NetworkPolicyList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "NetworkPolicy List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      namespace: $namespace,
      kubeGVRK: NetworkPolicyV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      networkPolicyListData = dl.data;
    }
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          name: item.metadata?.name,
          namespace: item.metadata?.namespace,
          kubeGVRK: NetworkPolicyV1GVRK,
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
    hrefRoot={routeString.networkPolicyList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Created At", ""]}
    items={networkPolicyListData?.items}
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
