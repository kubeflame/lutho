<script lang="ts">
  import { routeString, NodeV1GVRK } from "../lib/util";
  import type { V1NodeCondition, V1NodeList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import EmbeddedOptions from "../lib/tables/EmbeddedOptions.svelte";

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

  $: $dataSend = [
    {
      type: "list",
      request: {
        kubeGVRK: NodeV1GVRK,
      },
    },
  ];

  $: nodeListData = $dataList;

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

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar">
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.nodeList}
    isNamespaced={false}
    tableHead={["Name", "Created At", "Status", ""]}
    items={nodeListData?.items}
  >
    <!-- <div slot="nodeSchedule" let:item>
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
    </td> -->
  </ListTable>
</RouterPage>
