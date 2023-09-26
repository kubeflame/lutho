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

<RouterPage bind:error={$sockError} bind:loading={$isLoading}>
  <ResourceToolbar slot="resource-toolbar" />
  {#if listAllData}
    <Stats
      stats={[
        {
          title: "Nodes",
          value: listAllData[NodeV1GVRK.resource]?.items?.length ?? 0,
          iconType: "node",
          description: "status",
        },
        {
          title: "Pods",
          value: listAllData[PodV1GVRK.resource]?.items?.length ?? 0,
          iconType: "pod",
          description: `${3} in pending state`,
        },
        {
          title: "CronJobs",
          value: listAllData[CronJobV1GVRK.resource]?.items?.length ?? 0,
          iconType: "cronjob",
          description: `${3} in pending state`,
        },
        {
          title: "Jobs",
          value: listAllData[JobV1GVRK.resource]?.items?.length ?? 0,
          iconType: "job",
          description: `${3} in pending state`,
        },
        {
          title: "StatefulSets",
          value: listAllData[StatefulSetV1GVRK.resource]?.items?.length ?? 0,
          iconType: "statefulSet",
          description: `${3} in pending state`,
        },
        {
          title: "ReplicaSets",
          value: listAllData[ReplicaSetV1GVRK.resource]?.items?.length ?? 0,
          iconType: "hashtag",
          description: `${3} in pending state`,
        },
      ]}
    />
  {/if}
</RouterPage>
