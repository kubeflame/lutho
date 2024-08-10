<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    DeploymentV1GVRK,
    getLabels,
    jsonStringClassName,
    parseFieldSelector,
    EventV1GVRK,
    randomUUID,
  } from "../lib/util";
  import type { V1Deployment } from "@kubernetes/client-node";
  import HeaderElement from "../lib/HeaderElement.svelte";
  import Tabs from "../lib/Tabs.svelte";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import socketStore from "../lib/socketStore";
  import RouterPage from "../lib/RouterPage.svelte";
  import ResourceToolbarBreadcrumbs from "../lib/ResourceToolbarBreadcrumbs.svelte";
  import type { TabQueryParam } from "../lib/types";
  import Details from "../lib/Details.svelte";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";
  import CaptionTag from "../lib/CaptionTag.svelte";
  import Events from "../lib/Events.svelte";

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
  let deployData: V1Deployment;
  let codeMirrorChanged: boolean;
  let eventListData: any;

  $: toolbarContent = [
    { index: 0, name: "Deployment List", url: routeString.deployList },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: DeploymentV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendList = {
    type: "list",
    request: {
      name: params.name,
      kubeGVRK: EventV1GVRK,
      kubeOptions: {
        fieldSelector: parseFieldSelector({
          kind: DeploymentV1GVRK.kind,
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
      deployData = dg.data;
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
      deployData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: DeploymentV1GVRK,
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

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
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
          { name: "Name", value: deployData?.metadata?.name },
          { name: "Namespace", value: deployData?.metadata?.namespace },
          {
            name: "Creation Timestamp",
            value: deployData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: deployData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(deployData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(deployData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Replicas",
            value: deployData?.spec?.replicas,
          },
          {
            name: "Selector",
            value: JSON.stringify(deployData?.spec?.selector, null, 2) ?? "-",
            className: jsonStringClassName,
          },
          {
            name: "Strategy",
            value: JSON.stringify(deployData?.spec?.strategy, null, 2) ?? "-",
            className: jsonStringClassName,
          },
        ]}
      />
    </Details>

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Replicas",
            value: deployData?.status?.replicas,
          },
          {
            name: "Available Replicas",
            value: deployData?.status?.availableReplicas,
          },
          {
            name: "Updated Replicas",
            value: deployData?.status?.updatedReplicas,
          },
          {
            name: "Ready Replicas",
            value: deployData?.status?.readyReplicas,
          },
          {
            name: "Conditions",
            value:
              JSON.stringify(deployData?.status?.conditions, null, 2) ?? "-",
            className: jsonStringClassName,
          },
        ]}
      />
    </Details>

    <Details title={"Template"}>
      <EmbeddedTable
        tagName={"Spec"}
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Affinity",
            value:
              JSON.stringify(
                deployData?.spec?.template?.spec?.affinity,
                null,
                2,
              ) ?? "-",
            className: jsonStringClassName,
          },
          {
            name: "Automount Service Account Token",
            value:
              deployData?.spec?.template?.spec?.automountServiceAccountToken,
          },
          {
            name: "DNS Policy",
            value: deployData?.spec?.template?.spec?.dnsPolicy,
          },
          {
            name: "Restart Policy",
            value: deployData?.spec?.template?.spec?.restartPolicy,
          },
          {
            name: "Scheduler Name",
            value: deployData?.spec?.template?.spec?.schedulerName,
          },
          {
            name: "Security Context",
            value:
              JSON.stringify(
                deployData?.spec?.template?.spec?.securityContext,
                null,
                2,
              ) ?? "-",
            className: jsonStringClassName,
          },
          {
            name: "Service Account",
            value: deployData?.spec?.template?.spec?.serviceAccount,
          },
          {
            name: "Service Account Name",
            value: deployData?.spec?.template?.spec?.serviceAccountName,
          },
        ]}
      />

      {#if deployData?.spec?.template?.spec?.containers}
        <div class="card card-compact mt-2 rounded-lg">
          <CaptionTag tagName={"Containers"} />
          {#each deployData?.spec?.template?.spec?.containers as container}
            <div class="card-title pl-2 pt-2 text-sm font-light tracking-wider">
              <div class="badge badge-ghost badge-xs" />
              {container.name}
            </div>
            <div class="card-body space-y-3">
              <EmbeddedTable
                tableType={"custom-vertical"}
                tableItems={[
                  {
                    name: "Name",
                    value: container.name,
                  },
                  {
                    name: "Image",
                    value: container.image,
                  },
                  {
                    name: "Image Pull Policy",
                    value: container.imagePullPolicy,
                  },
                  {
                    name: "Ports",
                    value: JSON.stringify(container.ports, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                  {
                    name: "Env",
                    value: JSON.stringify(container.env, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                  {
                    name: "Security Context",
                    value:
                      JSON.stringify(container.securityContext, null, 2) ?? "-",
                    className:
                      "max-h-fit w-full rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
                  },
                ]}
              />
            </div>
          {/each}
        </div>
      {/if}
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(deployData)}
      bind:docStore
      bind:codeMirrorChanged
    />
  {:else if tabQueryParam === "events"}
    <Events kubeEvents={eventListData?.items} eventsReversed={0} />
  {/if}
</RouterPage>
