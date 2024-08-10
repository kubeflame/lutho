<script lang="ts">
  import CodeMirror from "../lib/codemirror/CodeMirror.svelte";
  import { dump as toYaml, load as fromYaml } from "js-yaml";
  import {
    tabs,
    routeString,
    ReplicationControllerV1GVRK,
    getLabels,
    randomUUID,
  } from "../lib/util";
  import type { V1ReplicationController } from "@kubernetes/client-node";
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

  export let params: any;

  const tabItems = [tabs.details, tabs.yaml];
  const { sockError, isLoading, dataSend, dataGet, dataUpdate } = socketStore();

  let tabQueryParam: TabQueryParam;
  let docStore: any;
  let replicationControllerData: V1ReplicationController;
  let codeMirrorChanged: boolean;

  $: toolbarContent = [
    {
      index: 0,
      name: "ReplicationController List",
      url: routeString.deployList,
    },
    { index: 1, name: params.name },
  ];

  $: sendGet = {
    opID: randomUUID(),
    type: "get",
    request: {
      namespace: params.namespace,
      name: params.name,
      kubeGVRK: ReplicationControllerV1GVRK,
    },
  } as any;

  $: $dataSend = [sendGet];

  $: sendUpdate = {
    opID: randomUUID(),
    type: "update",
  } as any;

  dataGet.subscribe((dg) => {
    if (dg && dg.op?.opID === sendGet.opID) {
      replicationControllerData = dg.data;
    }
  });

  dataUpdate.subscribe((du) => {
    if (du && du.op?.opID === sendUpdate.opID) {
      replicationControllerData = du.data;
    }
  });

  function update() {
    $dataSend = [
      {
        ...sendUpdate,
        request: {
          namespace: params.namespace,
          name: params.name,
          kubeGVRK: ReplicationControllerV1GVRK,
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
          { name: "Name", value: replicationControllerData?.metadata?.name },
          {
            name: "Namespace",
            value: replicationControllerData?.metadata?.namespace,
          },
          {
            name: "Creation Timestamp",
            value: replicationControllerData?.metadata?.creationTimestamp,
          },
          {
            name: "UID",
            value: replicationControllerData?.metadata?.uid,
          },
        ]}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Annotations"}
        tableItems={getLabels(replicationControllerData?.metadata?.annotations)}
      />

      <EmbeddedTable
        tableType={"badges"}
        tagName={"Labels"}
        tableItems={getLabels(replicationControllerData?.metadata?.labels)}
      />
    </Details>

    <Details title={"Spec"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "Selector",
            value:
              JSON.stringify(
                replicationControllerData?.spec?.selector,
                null,
                2,
              ) ?? "{}",
            className:
              "max-h-fit max-w-fit rounded-lg whitespace-pre border border-base-300 bg-base-200/10 px-1 font-mono text-sm",
          },
          {
            name: "replicas",
            value: replicationControllerData?.spec?.replicas,
          },
        ]}
      />

      {#if replicationControllerData?.spec?.template?.spec?.containers}
        <div class="card card-compact mt-2 rounded-lg">
          <CaptionTag tagName={"Containers"} />
          {#each replicationControllerData?.spec?.template?.spec?.containers as container}
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

    <Details title={"Status"}>
      <EmbeddedTable
        tableType={"custom-vertical"}
        tableItems={[
          {
            name: "availableReplicas",
            value: replicationControllerData?.status?.availableReplicas,
          },
          {
            name: "fullyLabeledReplicas",
            value: replicationControllerData?.status?.fullyLabeledReplicas,
          },
          {
            name: "replicas",
            value: replicationControllerData?.status?.replicas,
          },
          {
            name: "readyReplicas",
            value: replicationControllerData?.status?.readyReplicas,
          },
        ]}
      />
    </Details>
  {:else if tabQueryParam === "yaml"}
    <CodeMirror
      value={toYaml(replicationControllerData)}
      bind:codeMirrorChanged
      bind:docStore
    />
  {/if}
</RouterPage>
