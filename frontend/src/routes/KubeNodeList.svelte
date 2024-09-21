<script lang="ts">
  import { routeString, NodeV1GVRK, randomUUID } from "../lib/util";
  import type { V1NodeCondition, V1NodeList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";

  let nodeListData: V1NodeList;

  const getNodeStatus = (nc: V1NodeCondition[]): string => {
    return (
      nc.find(
        (item: V1NodeCondition): boolean =>
          item.type === "Ready" && item.status === "True",
      )?.type || "Not Ready"
    );
  };

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Node List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      kubeGVRK: NodeV1GVRK,
    },
  } as any;

  $: $dataSend = [sendList];

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      nodeListData = dl.data;
    }
  });

  // function onDelete(item: any) {
  //   $dataSend = [
  //     {
  //       type: "delete",
  //       request: {
  //         name: item.metadata?.name,
  //         kubeGVRK: NodeV1GVRK,
  //       },
  //     },
  //   ];
  // }
</script>

<HeaderElement />

<RouterPage bind:errorMessage={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.nodeList}
    isNamespaced={false}
    tableHead={["Name", "Created At", "Status", ""]}
    items={nodeListData?.items}
    embeddedOptionsDataFn={() => []}
    hasEmbeddedOptions={false}
  >
    <div slot="nodeSchedule" let:item>
      {#if item.spec?.unschedulable}
        <span
          class="badge badge-ghost badge-outline badge-sm font-light tracking-wide"
        >
          SchedulingDisabled
        </span>
      {/if}
    </div>
    <td slot="nodeStatus" let:item>
      {getNodeStatus(item.status?.conditions ?? [])}
    </td>
  </ListTable>
</RouterPage>
