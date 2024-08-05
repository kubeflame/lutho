<script lang="ts">
  import ResourceToolbar from "../ResourceToolbar.svelte";
  import RouterPage from "../RouterPage.svelte";
  import socketStore from "../socketStore";
  import { namespace } from "../stores";
  import type { DashboardList } from "../types";
  import {
    CronJobV1GVRK,
    DeploymentV1GVRK,
    JobV1GVRK,
    NodeV1GVRK,
    PodV1GVRK,
    ReplicaSetV1GVRK,
    ReplicationControllerV1GVRK,
    ServiceV1GVRK,
    StatefulSetV1GVRK,
  } from "../util";
  import Stats from "./Stats.svelte";

  let listAllData: DashboardList;

  const { sockError, isLoading, dataSend, dataList } = socketStore();

  $: $dataSend = [
    {
      type: "listAll",
      request: {
        namespace: $namespace,
        kubeGVRKList: [
          NodeV1GVRK,
          ReplicationControllerV1GVRK,
          PodV1GVRK,
          DeploymentV1GVRK,
          StatefulSetV1GVRK,
          JobV1GVRK,
          CronJobV1GVRK,
          ServiceV1GVRK,
          ReplicaSetV1GVRK,
        ],
      },
    },
  ];

  $: listAllData = $dataList;
</script>

<RouterPage bind:error={$sockError} loading={false}>
  <ResourceToolbar slot="resource-toolbar" />
  <div class="grid grid-flow-row grid-cols-3 gap-x-1 gap-y-2">
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "Nodes",
        value: listAllData
          ? listAllData[NodeV1GVRK.resource]?.items?.length
          : 0,
        iconType: "node",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "Deployments",
        value: listAllData
          ? listAllData[DeploymentV1GVRK.resource]?.items?.length
          : 0,
        iconType: "deployment",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "Pods",
        value: listAllData ? listAllData[PodV1GVRK.resource]?.items?.length : 0,
        iconType: "pod",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "CronJobs",
        value: listAllData
          ? listAllData[CronJobV1GVRK.resource]?.items?.length
          : 0,
        iconType: "cronjob",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "Jobs",
        value: listAllData ? listAllData[JobV1GVRK.resource]?.items?.length : 0,
        iconType: "job",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "StatefulSets",
        value: listAllData
          ? listAllData[StatefulSetV1GVRK.resource]?.items?.length
          : 0,
        iconType: "statefulSet",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "ReplicaSets",
        value: listAllData
          ? listAllData[ReplicaSetV1GVRK.resource]?.items?.length
          : 0,
        iconType: "hashtag",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "ReplicationController",
        value: listAllData
          ? listAllData[ReplicationControllerV1GVRK.resource]?.items?.length
          : 0,
        iconType: "hashtag",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={$isLoading}
      stat={{
        title: "Services",
        value: listAllData
          ? listAllData[ServiceV1GVRK.resource]?.items?.length
          : 0,
        iconType: "service",
        description: "",
      }}
    />
  </div>
</RouterPage>
