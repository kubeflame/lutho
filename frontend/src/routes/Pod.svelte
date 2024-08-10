<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import Shell from "../lib/xterm/Shell.svelte";
  import {
    routeString,
    tabs,
    EventV1GVRK,
    PodV1GVRK,
    getLabels,
    parseFieldSelector,
    randomUUID,
  } from "../lib/util";
  import { type EmbeddedTableItem, type TabQueryParam } from "../lib/types";
  import Events from "../lib/Events.svelte";
  import Details from "../lib/Details.svelte";
  import type {
    CoreV1EventList,
    V1Container,
    V1ContainerPort,
    V1Pod,
  } from "@kubernetes/client-node";
  import Logs from "../lib/Logs.svelte";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import { type Writable } from "svelte/store";

  export let params: any;

  const tabItems = [
    tabs.details,
    tabs.yaml,
    tabs.logs,
    tabs.shell,
    tabs.events,
  ];
  const {
    sockState,
    sockError,
    isLoading,
    dataSend,
    dataList,
    dataGet,
    dataUpdate,
  } = socketStore();

  let docStore: any;
  let followLogs: boolean;
  let onChangeLogsBtn: Writable<Event>;
  let podData: V1Pod;
  let tabQueryParam: TabQueryParam;
  let activeContainer: Writable<string>;
  let showShellReconnect: Writable<boolean>;
  let reconnectShell: boolean;
  let codeMirrorChanged: boolean;
  let eventListData: any;

  $: toolbarContent = [
    { index: 0, name: "Pod List", url: routeString.podList },
    { index: 1, name: params.name },
  ];
  $: containers = getContainers(podData?.spec?.containers ?? []);

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: PodV1GVRK,
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
          kind: PodV1GVRK.kind,
          name: params.name,
          namespace: params.namespace,
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
      podData = dg.data;
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
      podData = du.data;
    }
  });

  function getContainers(containers: V1Container[]): string[] {
    return containers.map((c) => {
      return c.name;
    });
  }

  function getContainerPorts(ports: V1ContainerPort[]): EmbeddedTableItem[] {
    let arr: EmbeddedTableItem[] = [];
    Object.entries(ports).forEach((port) => {
      Object.entries(port[1]).forEach((val) =>
        arr.push({ name: val[0].split(/(?=[A-Z])/).join(" "), value: val[1] }),
      );
    });
    return arr;
  }

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: PodV1GVRK,
          data: JSON.stringify(fromYaml($docStore)),
        },
      },
    ];
  }

  let unique = {};
  function restart() {
    unique = {};
    reconnectShell = false;
  }

  $: if (reconnectShell || $activeContainer) restart();
</script>

<HeaderElement>
  <Tabs
    slot="tabs"
    bind:activeContainer
    bind:tabQueryParam
    bind:containers
    tabQueryParamDefault={"details"}
    {tabItems}
  />
</HeaderElement>

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar
    slot="resource-toolbar"
    bind:codeMirrorChanged
    bind:followLogs
    bind:onChangeLogsBtn
    bind:showShellReconnect
    bind:reconnectShell
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
          { name: "Name", value: podData?.metadata?.name },
          { name: "Namespace", value: podData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: podData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: podData?.metadata?.uid,
          },
          {
            name: "Owner References",
            value:
              JSON.stringify(podData?.metadata?.ownerReferences, null, 2) ??
              "-",
            className:
              "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre",
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(podData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(podData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Affinity",
            value: JSON.stringify(podData?.spec?.affinity, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "DNS Policy",
            value: podData?.spec?.dnsPolicy,
          },
          {
            name: "Node Name",
            value: podData?.spec?.nodeName,
          },
          {
            name: "Restart Policy",
            value: podData?.spec?.restartPolicy,
          },
          {
            name: "Scheduler Name",
            value: podData?.spec?.schedulerName,
          },
          {
            name: "Security Context",
            value:
              JSON.stringify(podData?.spec?.securityContext, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Service Account",
            value: podData?.spec?.serviceAccount,
          },
          // {
          //   name: "Tolerations",
          //   value: JSON.stringify(podData?.spec?.tolerations, null, 2) ?? "{}",
          //   className:
          //     "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          // },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          // {
          //   name: "Conditions",
          //   value: podData?.status?.conditions,
          // },
          {
            name: "Container Statuses",
            value:
              JSON.stringify(podData?.status?.containerStatuses, null, 2) ??
              "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "Host IP",
            value: podData?.status?.hostIP,
          },
          {
            name: "Phase",
            value: podData?.status?.phase,
          },
          {
            name: "Pod IPs",
            value: JSON.stringify(podData?.status?.podIPs, null, 2) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "QOS Class",
            value: podData?.status?.qosClass,
          },
          {
            name: "Start Time",
            value: podData?.status?.startTime,
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
          { name: "Last Probe Time" },
          { name: "Last Transition Time" },
          { name: "Reason" },
          { name: "Message" },
        ]}
      >
        <tbody class="text-sm">
          {#if podData?.status?.conditions}
            {#each podData.status.conditions as condition}
              <tr>
                <td>{condition.type}</td>
                <td>{condition.status}</td>
                <td>{condition.lastProbeTime ?? "-"}</td>
                <td>{condition.lastTransitionTime ?? "-"}</td>
                <td>{condition.reason ?? "-"}</td>
                <td>{condition.message ?? "-"}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>

    <!-- <Details title={"Controlled By"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[
          { name: "Name" },
          { name: "Kind" },
          { name: "Controller" },
          { name: "Block Owner Deletion" },
        ]}
      >
        <tbody class="text-sm">
          {#if podData?.metadata?.ownerReferences}
            {#each podData?.metadata?.ownerReferences as ref}
              <tr>
                <td>{ref.name}</td>
                <td>{ref.kind}</td>
                <td>{ref.controller}</td>
                <td>{ref.blockOwnerDeletion}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details> -->

    <Details title={"Tolerations"}>
      <EmbeddedTable
        tableType={"custom-body"}
        tableItems={[{ name: "Key" }, { name: "Operator" }, { name: "Effect" }]}
      >
        <tbody class="text-sm">
          {#if podData?.spec?.tolerations}
            {#each podData?.spec?.tolerations as toleration}
              <tr>
                <td>{toleration.key ?? "-"}</td>
                <td>{toleration.operator ?? "-"}</td>
                <td>{toleration.effect ?? "-"}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </EmbeddedTable>
    </Details>

    <Details title={"Containers"}>
      {#if podData?.spec?.containers}
        {#each podData?.spec?.containers as container}
          <div class="card card-compact rounded-lg">
            <div class="card-title pl-2 pt-2 text-sm font-light tracking-wider">
              <div
                class="badge {podData?.status?.containerStatuses?.find(
                  (c) => c.name === container.name,
                )?.state?.running
                  ? 'badge-primary'
                  : 'badge-warning'} badge-xs"
              />
              {container.name}
            </div>
            <div class="card-body space-y-3">
              <EmbeddedTable
                tagName={"Summary"}
                tableType={"one_row"}
                tableItems={[
                  {
                    name: "Image",
                    value: container.image,
                  },
                  {
                    name: "Image Pull Policy",
                    value: container.imagePullPolicy,
                  },
                  {
                    name: "Ready",
                    value:
                      podData?.status?.containerStatuses?.find(
                        (c) => c.name === container.name,
                      )?.ready ?? "-",
                  },
                  {
                    name: "Started",
                    value:
                      podData?.status?.containerStatuses?.find(
                        (c) => c.name === container.name,
                      )?.started ?? "-",
                  },
                  {
                    name: "Restart Count",
                    value:
                      podData?.status?.containerStatuses?.find(
                        (c) => c.name === container.name,
                      )?.restartCount ?? "-",
                  },
                ]}
              />

              {#if container.env}
                <EmbeddedTable
                  tagName={"Environment Variables"}
                  tableType={"custom-body"}
                  tableItems={[
                    { name: "Name" },
                    { name: "Value" },
                    { name: "Value From" },
                  ]}
                >
                  <tbody class="">
                    {#each container.env as env}
                      <tr>
                        <td>{env.name}</td>
                        <td>{env.value ?? "-"}</td>
                        <td>
                          {#if env.valueFrom}
                            <div
                              class="border-base-300 bg-base-200/10 max-h-fit max-w-fit rounded-lg border px-1 font-mono text-sm"
                            >
                              {JSON.stringify(env.valueFrom, null, 2)}
                            </div>
                          {:else}
                            -
                          {/if}
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </EmbeddedTable>
              {/if}

              {#if container.ports}
                <EmbeddedTable
                  tagName={"Ports"}
                  tableType={"custom-body"}
                  tableItems={[
                    { name: "Name" },
                    { name: "Container Port" },
                    { name: "Host IP" },
                    { name: "Host Port" },
                    { name: "Protocol" },
                  ]}
                >
                  <tbody class="text-sm">
                    {#each container.ports as port}
                      <tr>
                        <td>{port.name}</td>
                        <td>{port.containerPort}</td>
                        <td>{port.hostIP ?? "-"}</td>
                        <td>{port.hostPort ?? "-"}</td>
                        <td>{port.protocol}</td>
                      </tr>
                    {/each}
                  </tbody>
                </EmbeddedTable>
              {/if}
            </div>
          </div>
        {/each}
      {/if}
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror value={toYaml(podData)} bind:codeMirrorChanged bind:docStore />
  {:else if tabQueryParam === "logs"}
    <Logs
      bind:followLogs
      bind:onChangeLogsBtn
      bind:activeContainer
      bind:name={params.name}
      bind:namespace={params.namespace}
    />
  {:else if tabQueryParam === "shell"}
    {#key unique}
      <Shell
        bind:name={params.name}
        bind:namespace={params.namespace}
        bind:showShellReconnect
        bind:activeContainer
      />
    {/key}
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
