<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    EventV1GVRK,
    NodeV1GVRK,
    getLabels,
    parseFieldSelector,
    randomUUID,
  } from "../lib/util";
  import Events from "../lib/Events.svelte";
  import type { V1Node } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Details from "../lib/Details.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml, tabs.events];
  const {
    sockState,
    sockError,
    isLoading,
    dataSend,
    dataList,
    dataGet,
    dataUpdate,
  } = socketStore();

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let nodeData: V1Node;
  let codeMirrorChanged: boolean;
  let eventListData: any;

  $: toolbarContent = [
    { index: 0, name: "Node List", url: routeString.nodeList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      name: params.name,
      kubeGVRK: NodeV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendList = {
    opID: randomUUID(),
    type: "list",
    request: {
      name: params.name,
      kubeGVRK: EventV1GVRK,
      kubeOptions: {
        fieldSelector: parseFieldSelector({
          kind: NodeV1GVRK.kind,
          name: params.name,
          namespace: "",
        }),
      },
    },
  } as any;

  $: if (tabQueryParam === "events") {
    $dataSend = [sendList];
  }

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      nodeData = dg.data;
    }
  });

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.opID) {
      eventListData = dl.data;
    }
  });

  dataSend.subscribe((ds) => {
    if (ds && $sockState.state === WebSocket.CLOSED) $sockState.refresh = true;
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      nodeData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          name: params.name,
          kubeGVRK: NodeV1GVRK,
          data: JSON.stringify(fromYaml($docStore)),
        },
      },
    ];
  }
</script>

<HeaderElement>
  <Tabs
    slot="tabs"
    bind:tabQueryParam
    tabQueryParamDefault={"details"}
    {tabItems}
  />
</HeaderElement>

<RouterPage bind:errorMessage={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar
    slot="resource-toolbar"
    bind:codeMirrorChanged
    bind:tabQueryParam
    onClickSubmit={update}
  >
    <ResourceToolbarBreadcrumbs slot="breadcrumbs" bind:toolbarContent />
  </ResourceToolbar>
  {#if tabQueryParam === "details"}
    <Details title={"Resource Information"}>
      <EmbeddedTable
        tagName={"Summary"}
        tableType={"custom-vertical"}
        tableItems={[
          { name: "Name", value: nodeData?.metadata?.name },
          {
            name: "Creation Timestamp",
            value: nodeData?.metadata?.creationTimestamp,
          },
          {
            name: "Finalizers",
            value:
              JSON.stringify(nodeData?.metadata?.finalizers, null, 2) ?? "-",
            className:
              "max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "UID",
            value: nodeData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(nodeData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(nodeData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "ProviderID",
            value: nodeData?.spec?.providerID,
          },
          { name: "PodCIDR", value: nodeData?.spec?.podCIDR },
          {
            name: "PodCIDRs",
            value: JSON.stringify(nodeData?.spec?.podCIDRs, null, 2) ?? "-",
            className:
              "max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
        ]}
      />
    </Details>

    <Details title={"Node Info"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Architecture",
            value: nodeData?.status?.nodeInfo?.architecture,
          },
          {
            name: "BootID",
            value: nodeData?.status?.nodeInfo?.bootID,
          },
          {
            name: "ContainerRuntimeVersion",
            value: nodeData?.status?.nodeInfo?.containerRuntimeVersion,
          },
          {
            name: "KernelVersion",
            value: nodeData?.status?.nodeInfo?.kernelVersion,
          },
          {
            name: "KubeProxyVersion",
            value: nodeData?.status?.nodeInfo?.kubeProxyVersion,
          },
          {
            name: "KubeletVersion",
            value: nodeData?.status?.nodeInfo?.kubeletVersion,
          },
          {
            name: "MachineID",
            value: nodeData?.status?.nodeInfo?.machineID,
          },
          {
            name: "OperatingSystem",
            value: nodeData?.status?.nodeInfo?.operatingSystem,
          },
          {
            name: "OSImage",
            value: nodeData?.status?.nodeInfo?.osImage,
          },
          {
            name: "SystemUUID",
            value: nodeData?.status?.nodeInfo?.systemUUID,
          },
        ]}
      />
    </Details>

    <Details title={"Conditions"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[
          { name: "Type" },
          { name: "Status" },
          { name: "Last Heartbeat Time" },
          { name: "Last Transition Time" },
          { name: "Reason" },
          { name: "Message" },
        ]}
      >
        <tbody class="text-sm">
          {#if nodeData?.status?.conditions}
            {#each nodeData.status.conditions as condition}
              <tr>
                <td>{condition.type}</td>
                <td>{condition.status}</td>
                <td>{condition.lastHeartbeatTime ?? "-"}</td>
                <td>{condition.lastTransitionTime ?? "-"}</td>
                <td>{condition.reason ?? "-"}</td>
                <td>{condition.message ?? "-"}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>

    <Details title={"Addresses"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[{ name: "Type" }, { name: "Address" }]}
      >
        <tbody class="text-sm">
          {#if nodeData?.status?.addresses}
            {#each nodeData?.status?.addresses as address}
              <tr>
                <td>{address.type ?? "-"}</td>
                <td>{address.address ?? "-"}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>

    <Details title={"Allocatable"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "CPU",
            value: nodeData?.status?.allocatable?.cpu,
          },
          {
            name: "Memory",
            value: nodeData?.status?.allocatable?.memory,
          },
          {
            name: "Pods",
            value: nodeData?.status?.allocatable?.pods,
          },
          {
            name: "Ephemeral Storage",
            value: nodeData?.status?.allocatable?.["ephemeral-storage"],
          },
          {
            name: "Hugepages 1Gi",
            value: nodeData?.status?.allocatable?.["hugepages-1Gi"],
          },
          {
            name: "Hugepages 32Mi",
            value: nodeData?.status?.allocatable?.["hugepages-32Mi"],
          },
          {
            name: "Hugepages 2Mi",
            value: nodeData?.status?.allocatable?.["hugepages-2Mi"],
          },
          {
            name: "Hugepages 64Ki",
            value: nodeData?.status?.allocatable?.["hugepages-64Ki"],
          },
        ]}
      />
    </Details>

    <Details title={"Capacity"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "CPU",
            value: nodeData?.status?.capacity?.cpu,
          },
          {
            name: "Memory",
            value: nodeData?.status?.capacity?.memory,
          },
          {
            name: "Pods",
            value: nodeData?.status?.capacity?.pods,
          },
          {
            name: "Ephemeral Storage",
            value: nodeData?.status?.capacity?.["ephemeral-storage"],
          },
          {
            name: "Hugepages 1Gi",
            value: nodeData?.status?.capacity?.["hugepages-1Gi"],
          },
          {
            name: "Hugepages 32Mi",
            value: nodeData?.status?.capacity?.["hugepages-32Mi"],
          },
          {
            name: "Hugepages 2Mi",
            value: nodeData?.status?.capacity?.["hugepages-2Mi"],
          },
          {
            name: "Hugepages 64Ki",
            value: nodeData?.status?.capacity?.["hugepages-64Ki"],
          },
        ]}
      />
    </Details>

    <Details title={"Images"}></Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror value={toYaml(nodeData)} bind:codeMirrorChanged bind:docStore />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
