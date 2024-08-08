import HomePage from "./routes/HomePage.svelte";
import KubeNode from "./routes/KubeNode.svelte";
import KubeNodeList from "./routes/KubeNodeList.svelte";
import PodList from "./routes/PodList.svelte";
import Pod from "./routes/Pod.svelte";
import DeploymentList from "./routes/DeploymentList.svelte";
import Deployment from "./routes/Deployment.svelte";
import SecretList from "./routes/SecretList.svelte";
import Secret from "./routes/Secret.svelte";
import ConfigMapList from "./routes/ConfigMapList.svelte";
import ConfigMap from "./routes/ConfigMap.svelte";
import StatefulSetList from "./routes/StatefulSetList.svelte";
import StatefulSet from "./routes/StatefulSet.svelte";
import JobList from "./routes/JobList.svelte";
import Job from "./routes/Job.svelte";
import IngressList from "./routes/IngressList.svelte";
import Ingress from "./routes/Ingress.svelte";
import NotFound from "./lib/NotFound.svelte";
import ServiceList from "./routes/ServiceList.svelte";
import Service from "./routes/Service.svelte";
import ClusterRole from "./routes/ClusterRole.svelte";
import ClusterRoleList from "./routes/ClusterRoleList.svelte";
import ClusterRoleBinding from "./routes/ClusterRoleBinding.svelte";
import ClusterRoleBindingList from "./routes/ClusterRoleBindingList.svelte";
import Role from "./routes/Role.svelte";
import RoleList from "./routes/RoleList.svelte";
import RoleBinding from "./routes/RoleBinding.svelte";
import RoleBindingList from "./routes/RoleBindingList.svelte";
import CronJob from "./routes/CronJob.svelte";
import CronJobList from "./routes/CronJobList.svelte";
import DaemonSet from "./routes/DaemonSet.svelte";
import DaemonSetList from "./routes/DaemonSetList.svelte";
import EventList from "./routes/EventList.svelte";
import IngressClass from "./routes/IngressClass.svelte";
import IngressClassList from "./routes/IngressClassList.svelte";
import Namespace from "./routes/Namespace.svelte";
import NamespaceList from "./routes/NamespaceList.svelte";
import NetworkPolicy from "./routes/NetworkPolicy.svelte";
import NetworkPolicyList from "./routes/NetworkPolicyList.svelte";
import PersistentVolume from "./routes/PersistentVolume.svelte";
import PersistentVolumeList from "./routes/PersistentVolumeList.svelte";
import PersistentVolumeClaim from "./routes/PersistentVolumeClaim.svelte";
import PersistentVolumeClaimList from "./routes/PersistentVolumeClaimList.svelte";
import ReplicationController from "./routes/ReplicationController.svelte";
import ReplicationControllerList from "./routes/ReplicationControllerList.svelte";
import ServiceAccount from "./routes/ServiceAccount.svelte";
import ServiceAccountList from "./routes/ServiceAccountList.svelte";
import StorageClass from "./routes/StorageClass.svelte";
import StorageClassList from "./routes/StorageClassList.svelte";
import Event from "./routes/Event.svelte";
import ReplicaSetList from "./routes/ReplicaSetList.svelte";
import ReplicaSet from "./routes/ReplicaSet.svelte";
import HelmHome from "./routes/helm/HelmHome.svelte";
import ListReleases from "./routes/helm/ListReleases.svelte";
import ViewRelease from "./routes/helm/ViewRelease.svelte";
import ListRepos from "./routes/helm/ListRepos.svelte";
import ViewRepo from "./routes/helm/ViewRepo.svelte";
import AuthPage from "./lib/AuthPage.svelte";

export default {
  "/": HomePage,

  "/auth": AuthPage,

  "/nodes": KubeNodeList,
  "/nodes/:name": KubeNode,

  "/pods": PodList,
  "/pods/:namespace/:name": Pod,

  "/deployments": DeploymentList,
  "/deployments/:namespace/:name": Deployment,

  "/services": ServiceList,
  "/services/:namespace/:name": Service,

  "/secrets": SecretList,
  "/secrets/:namespace/:name": Secret,

  "/configmaps": ConfigMapList,
  "/configmaps/:namespace/:name": ConfigMap,

  "/statefulsets": StatefulSetList,
  "/statefulsets/:namespace/:name": StatefulSet,

  "/jobs": JobList,
  "/jobs/:namespace/:name": Job,

  "/ingresses": IngressList,
  "/ingresses/:namespace/:name": Ingress,

  "/clusterroles": ClusterRoleList,
  "/clusterroles/:name": ClusterRole,

  "/clusterrolebindings": ClusterRoleBindingList,
  "/clusterrolebindings/:name": ClusterRoleBinding,

  "/roles": RoleList,
  "/roles/:namespace/:name": Role,

  "/rolebindings": RoleBindingList,
  "/rolebindings/:namespace/:name": RoleBinding,

  "/cronjobs": CronJobList,
  "/cronjobs/:namespace/:name": CronJob,

  "/daemonsets": DaemonSetList,
  "/daemonsets/:namespace/:name": DaemonSet,

  "/replicasets": ReplicaSetList,
  "/replicasets/:namespace/:name": ReplicaSet,

  "/events": EventList,
  "/events/:namespace/:name": Event,

  "/ingressclasses": IngressClassList,
  "/ingressclasses/:name": IngressClass,

  "/namespaces": NamespaceList,
  "/namespaces/:name": Namespace,

  "/networkpolicies": NetworkPolicyList,
  "/networkpolicies/:namespace/:name": NetworkPolicy,

  "/persistentvolumes": PersistentVolumeList,
  "/persistentvolumes/:name": PersistentVolume,

  "/persistentvolumeclaims": PersistentVolumeClaimList,
  "/persistentvolumeclaims/:namespace/:name": PersistentVolumeClaim,

  "/replicationcontrollers": ReplicationControllerList,
  "/replicationcontrollesr/:namespace/:name": ReplicationController,

  "/serviceaccounts": ServiceAccountList,
  "/serviceaccounts/:namespace/:name": ServiceAccount,

  "/storageclasses": StorageClassList,
  "/storageclasses/:name": StorageClass,

  "/helm": HelmHome,
  "/helm/list": ListReleases,
  "/helm/releases/:namespace/:name": ViewRelease,
  "/helm/repository": ListRepos,
  "/helm/repository/:name": ViewRepo,

  "*": NotFound,
};
