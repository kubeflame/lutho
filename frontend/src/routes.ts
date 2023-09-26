import HomePage from "./routes/HomePage.svelte";
import About from "./routes/About.svelte";
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
import HelmList from "./routes/HelmList.svelte";
import HelmRelease from "./routes/HelmRelease.svelte";
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

export default {
  "/": HomePage,

  "/about": About,
  "/about/*": About,

  "/node": KubeNodeList,
  "/node/:name": KubeNode,

  "/pod": PodList,
  "/pod/:namespace/:name": Pod,

  "/deployment": DeploymentList,
  "/deployment/:namespace/:name": Deployment,

  "/service": ServiceList,
  "/service/:namespace/:name": Service,

  "/secret": SecretList,
  "/secret/:namespace/:name": Secret,

  "/configmap": ConfigMapList,
  "/configmap/:namespace/:name": ConfigMap,

  "/statefulset": StatefulSetList,
  "/statefulset/:namespace/:name": StatefulSet,

  "/job": JobList,
  "/job/:namespace/:name": Job,

  "/ingress": IngressList,
  "/ingress/:namespace/:name": Ingress,

  "/clusterrole": ClusterRoleList,
  "/clusterrole/:name": ClusterRole,

  "/clusterrolebinding": ClusterRoleBindingList,
  "/clusterrolebinding/:name": ClusterRoleBinding,

  "/role": RoleList,
  "/role/:namespace/:name": Role,

  "/rolebinding": RoleBindingList,
  "/rolebinding/:namespace/:name": RoleBinding,

  "/cronjob": CronJobList,
  "/cronjob/:namespace/:name": CronJob,

  "/daemonset": DaemonSetList,
  "/daemonset/:namespace/:name": DaemonSet,

  "/event": EventList,
  "/event/:namespace/:name": Event,

  "/ingressclass": IngressClassList,
  "/ingressclass/:name": IngressClass,

  "/namespace": NamespaceList,
  "/namespace/:name": Namespace,

  "/networkpolicy": NetworkPolicyList,
  "/networkpolicy/:namespace/:name": NetworkPolicy,

  "/persistentvolume": PersistentVolumeList,
  "/persistentvolume/:name": PersistentVolume,

  "/persistentvolumeclaim": PersistentVolumeClaimList,
  "/persistentvolumeclaim/:namespace/:name": PersistentVolumeClaim,

  "/replicationcontroller": ReplicationControllerList,
  "/replicationcontroller/:namespace/:name": ReplicationController,

  "/serviceaccount": ServiceAccountList,
  "/serviceaccount/:namespace/:name": ServiceAccount,

  "/storageclass": StorageClassList,
  "/storageclass/:name": StorageClass,

  "/release": HelmList,
  "/release/:namespace/:name": HelmRelease,

  "*": NotFound,
};
