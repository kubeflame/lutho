import HomePage from "./routes/HomePage.svelte";
import NotFound from "./lib/NotFound.svelte";
import AuthPage from "./lib/AuthPage.svelte";

import { wrap } from "svelte-spa-router/wrap";

export default {
  "/": HomePage,

  "/auth": AuthPage,

  "/nodes": wrap({
    asyncComponent: () => import("./routes/KubeNodeList.svelte"),
  }),
  "/nodes/:name": wrap({
    asyncComponent: () => import("./routes/KubeNode.svelte"),
  }),

  "/pods": wrap({
    asyncComponent: () => import("./routes/PodList.svelte"),
  }),
  "/pods/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Pod.svelte"),
  }),

  "/deployments": wrap({
    asyncComponent: () => import("./routes/DeploymentList.svelte"),
  }),
  "/deployments/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Deployment.svelte"),
  }),

  "/services": wrap({
    asyncComponent: () => import("./routes/ServiceList.svelte"),
  }),
  "/services/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Service.svelte"),
  }),

  "/secrets": wrap({
    asyncComponent: () => import("./routes/SecretList.svelte"),
  }),
  "/secrets/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Secret.svelte"),
  }),

  "/configmaps": wrap({
    asyncComponent: () => import("./routes/ConfigMapList.svelte"),
  }),
  "/configmaps/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/ConfigMap.svelte"),
  }),

  "/statefulsets": wrap({
    asyncComponent: () => import("./routes/StatefulSetList.svelte"),
  }),
  "/statefulsets/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/StatefulSet.svelte"),
  }),

  "/jobs": wrap({
    asyncComponent: () => import("./routes/JobList.svelte"),
  }),
  "/jobs/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Job.svelte"),
  }),

  "/ingresses": wrap({
    asyncComponent: () => import("./routes/IngressList.svelte"),
  }),
  "/ingresses/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Ingress.svelte"),
  }),

  "/clusterroles": wrap({
    asyncComponent: () => import("./routes/ClusterRoleList.svelte"),
  }),
  "/clusterroles/:name": wrap({
    asyncComponent: () => import("./routes/ClusterRole.svelte"),
  }),

  "/clusterrolebindings": wrap({
    asyncComponent: () => import("./routes/ClusterRoleBindingList.svelte"),
  }),
  "/clusterrolebindings/:name": wrap({
    asyncComponent: () => import("./routes/ClusterRoleBinding.svelte"),
  }),

  "/roles": wrap({
    asyncComponent: () => import("./routes/RoleList.svelte"),
  }),
  "/roles/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Role.svelte"),
  }),

  "/rolebindings": wrap({
    asyncComponent: () => import("./routes/RoleBindingList.svelte"),
  }),
  "/rolebindings/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/RoleBinding.svelte"),
  }),

  "/cronjobs": wrap({
    asyncComponent: () => import("./routes/CronJobList.svelte"),
  }),
  "/cronjobs/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/CronJob.svelte"),
  }),

  "/daemonsets": wrap({
    asyncComponent: () => import("./routes/DaemonSetList.svelte"),
  }),
  "/daemonsets/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/DaemonSet.svelte"),
  }),

  "/replicasets": wrap({
    asyncComponent: () => import("./routes/ReplicaSetList.svelte"),
  }),
  "/replicasets/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/ReplicaSet.svelte"),
  }),

  "/events": wrap({
    asyncComponent: () => import("./routes/EventList.svelte"),
  }),
  "/events/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/Event.svelte"),
  }),

  "/ingressclasses": wrap({
    asyncComponent: () => import("./routes/IngressClassList.svelte"),
  }),
  "/ingressclasses/:name": wrap({
    asyncComponent: () => import("./routes/IngressClass.svelte"),
  }),

  "/namespaces": wrap({
    asyncComponent: () => import("./routes/NamespaceList.svelte"),
  }),
  "/namespaces/:name": wrap({
    asyncComponent: () => import("./routes/Namespace.svelte"),
  }),

  "/networkpolicies": wrap({
    asyncComponent: () => import("./routes/NetworkPolicyList.svelte"),
  }),
  "/networkpolicies/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/NetworkPolicy.svelte"),
  }),

  "/persistentvolumes": wrap({
    asyncComponent: () => import("./routes/PersistentVolumeList.svelte"),
  }),
  "/persistentvolumes/:name": wrap({
    asyncComponent: () => import("./routes/PersistentVolume.svelte"),
  }),

  "/persistentvolumeclaims": wrap({
    asyncComponent: () => import("./routes/PersistentVolumeClaimList.svelte"),
  }),
  "/persistentvolumeclaims/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/PersistentVolumeClaim.svelte"),
  }),

  "/replicationcontrollers": wrap({
    asyncComponent: () => import("./routes/ReplicationControllerList.svelte"),
  }),
  "/replicationcontrollesr/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/ReplicationController.svelte"),
  }),

  "/serviceaccounts": wrap({
    asyncComponent: () => import("./routes/ServiceAccountList.svelte"),
  }),
  "/serviceaccounts/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/ServiceAccount.svelte"),
  }),

  "/storageclasses": wrap({
    asyncComponent: () => import("./routes/StorageClassList.svelte"),
  }),
  "/storageclasses/:name": wrap({
    asyncComponent: () => import("./routes/StorageClass.svelte"),
  }),

  "/helm": wrap({
    asyncComponent: () => import("./routes/helm/HelmHome.svelte"),
  }),
  "/helm/list": wrap({
    asyncComponent: () => import("./routes/helm/ListReleases.svelte"),
  }),
  "/helm/releases/:namespace/:name": wrap({
    asyncComponent: () => import("./routes/helm/ViewRelease.svelte"),
  }),
  "/helm/repository": wrap({
    asyncComponent: () => import("./routes/helm/ListRepos.svelte"),
  }),
  "/helm/repository/:name": wrap({
    asyncComponent: () => import("./routes/helm/ViewRepo.svelte"),
  }),

  "*": NotFound,
};
