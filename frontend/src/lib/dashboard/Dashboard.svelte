<script lang="ts">
  import ResourceToolbar from "../ResourceToolbar.svelte";
  import RouterPage from "../RouterPage.svelte";
  import socketStore from "../socketStore";
  import { namespace } from "../stores";
  import type { KubeOp } from "../types";
  import {
    CronJobV1GVRK,
    DeploymentV1GVRK,
    JobV1GVRK,
    NodeV1GVRK,
    PodV1GVRK,
    randomUUID,
    ReplicaSetV1GVRK,
    ReplicationControllerV1GVRK,
    ServiceV1GVRK,
    StatefulSetV1GVRK,
  } from "../util";
  import Stats from "./Stats.svelte";

  let nodeList: any;
  let deployList: any;
  let podList: any;
  let cronJobList: any;
  let jobList: any;
  let statefulSetList: any;
  let replicaSetList: any;
  let replicationControllerList: any;
  let serviceList: any;

  let nodeLoading: boolean = false;
  let deployLoading: boolean = false;
  let podLoading: boolean = false;
  let cronJobLoading: boolean = false;
  let jobLoading: boolean = false;
  let statefulSetLoading: boolean = false;
  let replicaSetLoading: boolean = false;
  let replicationControllerLoading: boolean = false;
  let serviceLoading: boolean = false;

  const { sockError, isLoading, dataSend, dataList } = socketStore();

  $: sendList = {
    nodes: {
      opID: randomUUID(),
      type: "list",
      request: {
        kubeGVRK: NodeV1GVRK,
      },
    },
    deployments: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: DeploymentV1GVRK,
      },
    },
    pods: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: PodV1GVRK,
      },
    },
    cronjobs: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: CronJobV1GVRK,
      },
    },
    jobs: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: JobV1GVRK,
      },
    },
    statefulsets: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: StatefulSetV1GVRK,
      },
    },
    replicasets: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: ReplicaSetV1GVRK,
      },
    },
    replicationcontrollers: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: ReplicationControllerV1GVRK,
      },
    },
    services: {
      opID: randomUUID(),
      type: "list",
      request: {
        namespace: $namespace,
        kubeGVRK: ServiceV1GVRK,
      },
    },
  } as Record<string, KubeOp>;

  $: $dataSend = [
    sendList.nodes,
    sendList.deployments,
    sendList.pods,
    sendList.cronjobs,
    sendList.jobs,
    sendList.statefulsets,
    sendList.replicasets,
    sendList.replicationcontrollers,
    sendList.services,
  ];

  dataSend.subscribe((ds) => {
    ds &&
      ds.forEach((d) => {
        if (d.opID === sendList.nodes?.opID) {
          nodeLoading = true;
        } else if (d.opID === sendList.deployments?.opID) {
          deployLoading = true;
        } else if (d.opID === sendList.pods?.opID) {
          podLoading = true;
        } else if (d.opID === sendList.cronjobs?.opID) {
          cronJobLoading = true;
        } else if (d.opID === sendList.jobs?.opID) {
          jobLoading = true;
        } else if (d.opID === sendList.statefulsets?.opID) {
          statefulSetLoading = true;
        } else if (d.opID === sendList.replicasets?.opID) {
          replicaSetLoading = true;
        } else if (d.opID === sendList.replicationcontrollers?.opID) {
          replicationControllerLoading = true;
        } else if (d.opID === sendList.services?.opID) {
          serviceLoading = true;
        }
      });
  });

  dataList.subscribe((dl) => {
    if (dl && dl.op?.opID === sendList.nodes?.opID) {
      nodeLoading = false;
      nodeList = dl.data;
    } else if (dl && dl.op?.opID === sendList.deployments?.opID) {
      deployLoading = false;
      deployList = dl.data;
    } else if (dl && dl.op?.opID === sendList.pods?.opID) {
      podLoading = false;
      podList = dl.data;
    } else if (dl && dl.op?.opID === sendList.cronjobs?.opID) {
      cronJobLoading = false;
      cronJobList = dl.data;
    } else if (dl && dl.op?.opID === sendList.jobs?.opID) {
      jobLoading = false;
      jobList = dl.data;
    } else if (dl && dl.op?.opID === sendList.statefulsets?.opID) {
      statefulSetLoading = false;
      statefulSetList = dl.data;
    } else if (dl && dl.op?.opID === sendList.replicasets?.opID) {
      replicaSetLoading = false;
      replicaSetList = dl.data;
    } else if (dl && dl.op?.opID === sendList.replicationcontrollers?.opID) {
      replicationControllerLoading = false;
      replicationControllerList = dl.data;
    } else if (dl && dl.op?.opID === sendList.services?.opID) {
      serviceLoading = false;
      serviceList = dl.data;
    }
  });
</script>

<RouterPage bind:error={$sockError} loading={false}>
  <ResourceToolbar slot="resource-toolbar" />
  <div class="grid grid-flow-row grid-cols-3 gap-x-1 gap-y-2">
    <Stats
      bind:isLoading={nodeLoading}
      stat={{
        title: "Nodes",
        value: nodeList ? nodeList?.items?.length : 0,
        iconType: "node",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={deployLoading}
      stat={{
        title: "Deployments",
        value: deployList ? deployList.items?.length : 0,
        iconType: "deployment",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={podLoading}
      stat={{
        title: "Pods",
        value: podList ? podList?.items?.length : 0,
        iconType: "pod",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={cronJobLoading}
      stat={{
        title: "CronJobs",
        value: cronJobList ? cronJobList?.items?.length : 0,
        iconType: "cronjob",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={jobLoading}
      stat={{
        title: "Jobs",
        value: jobList ? jobList?.items?.length : 0,
        iconType: "job",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={statefulSetLoading}
      stat={{
        title: "StatefulSets",
        value: statefulSetList ? statefulSetList?.items?.length : 0,
        iconType: "statefulSet",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={replicaSetLoading}
      stat={{
        title: "ReplicaSets",
        value: replicaSetList ? replicaSetList?.items?.length : 0,
        iconType: "hashtag",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={replicationControllerLoading}
      stat={{
        title: "ReplicationController",
        value: replicationControllerList
          ? replicationControllerList?.items?.length
          : 0,
        iconType: "hashtag",
        description: "",
      }}
    />
    <Stats
      bind:isLoading={serviceLoading}
      stat={{
        title: "Services",
        value: serviceList ? serviceList?.items?.length : 0,
        iconType: "service",
        description: "",
      }}
    />
  </div>
</RouterPage>
