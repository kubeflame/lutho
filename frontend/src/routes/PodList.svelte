<script lang="ts">
  import { namespace } from "../lib/stores";
  import {
    PodV1GVRK,
    getContainers,
    randomUUID,
    routeString,
  } from "../lib/util";
  import type { V1PodList } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import ListTable from "../lib/tables/ListTable.svelte";
  import { push } from "svelte-spa-router";

  let podListData: V1PodList;

  const { sockError, isLoading, dataSend, dataList, dataDelete } =
    socketStore();

  $: toolbarContent = [{ index: 0, name: "Pod List" }];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      namespace: $namespace,
      kubeGVRK: PodV1GVRK,
      // kubeOptions: {
      //   limit: 20,
      // },
    },
  } as any;

  $: $dataSend = [sendList];

  // $: kubeContinue = (podListData?.metadata as V1KubeListMeta).continue;

  // $: console.log(
  //   "continue - ",
  //   (podListData?.metadata as V1KubeListMeta).continue,
  // );
  // $: console.log(
  //   "remainingItemCount - ",
  //   podListData?.metadata?.remainingItemCount,
  // );

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      podListData = dl.data;
    }
  });

  function onDelete(item: any) {
    $dataSend = [
      {
        type: "delete",
        request: {
          namespace: item.metadata?.namespace,
          name: item.metadata?.name,
          kubeGVRK: PodV1GVRK,
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
    <!-- <button
      slot="custom"
      class="btn btn-sm"
      disabled={kubeContinue ? false : true}
      on:click={() => {
        $dataSend = [
          {
            type: "list",
            request: {
              namespace: $namespace,
              kubeGVRK: PodV1GVRK,
              kubeOptions: {
                limit: 20,
                continue: kubeContinue,
              },
            },
          },
        ];
      }}
    >
      Next
    </button> -->
  </ResourceToolbar>

  <ListTable
    hrefRoot={routeString.podList}
    isNamespaced={true}
    tableHead={["Name", "Namespace", "Status", "Created At", ""]}
    items={podListData?.items}
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
      {
        fn: () => {
          push(
            `${routeString.podList}/${item.metadata?.namespace}/${item.metadata?.name}?tab=shell`,
          );
        },
        classes: "hover:btn-warning",
        icon: "cmdLine",
        url: `${routeString.podList}/${item.metadata?.namespace}/${item.metadata?.name}`,
        hide: item.status?.phase !== "Running",
      },
    ]}
  >
    <div
      slot="podPhaseBadge"
      let:item
      class="badge badge-xs {item.status?.phase === 'Running'
        ? 'badge-primary'
        : 'badge-warning'}"
    />

    <td slot="podPhase" let:item>{item.status?.phase}</td>
  </ListTable>
</RouterPage>
