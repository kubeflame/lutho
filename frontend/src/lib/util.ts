import type { V1Container } from "@kubernetes/client-node";
import type {
  TabContent,
  KubeGVRK,
  HTTPStatusCode,
  EmbeddedTableItem,
  KubeFieldSelector,
} from "./types";

export const httpStatus: HTTPStatusCode = {
  200: "OK",
  201: "Created",
  400: "Bad Request",
  401: "Unauthorized",
  403: "Forbidden",
  404: "Not Found",
  405: "Method Not Allowed",
  408: "Request Timeout",
  418: "I'm a teapot",
  500: "Internal Server Error",
  501: "Not Implemented",
  502: "Bad Gateway",
  503: "Service Unavailable",
  504: "Gateway Timeout",
};

export const apiURL = {
  auth: "/srv/auth",
  authState: "/srv/auth/state",
  data: "/srv/data",
  dataWS: "/srv/data/ws",
  logs: "/srv/logs",
  logsStream: "/srv/logs/stream",
  shell: "/srv/shell",
  shellExec: "/srv/shell/exec",
} as const;

export const colorScheme = {
  light: {
    primary: "#65cc89",
    "primary-content": "#233e31",
    secondary: "#367bfb",
    "secondary-content": "#f9fafb",
    accent: "#e95234",
    "accent-content": "#f9fafb",
    neutral: "#333c4c",
    "neutral-content": "#f9fafb",
    "base-100": "#e9e7e7",
    "base-200": "#d7d5d5",
    "base-300": "#c7c2c2",
    "base-content": "#333c4c",
    info: "#3abef7",
    "info-content": "#012b3e",
    success: "#37d39a",
    "success-content": "#013321",
    warning: "#fabd22",
    "warning-content": "#382800",
    error: "#f87272",
    "error-content": "#470000",
  },
  dark: {
    primary: "#065f46",
    "primary-content": "#ffffff",
    secondary: "#0ea5e9",
    "secondary-content": "#bae6fd",
    accent: "#fb7185",
    "accent-content": "#ffffff",
    neutral: "#414557",
    "neutral-content": "#d6d7db",
    "base-100": "#262935",
    "base-200": "#181921",
    "base-300": "#09090d",
    "base-content": "#f8f8f1",
    info: "#8be8fd",
    "info-content": "#202e31",
    success: "#52f97c",
    "success-content": "#192e1c",
    warning: "#db7706",
    "warning-content": "#2d2e1e",
    error: "#ff5756",
    "error-content": "#311816",
  },
} as const;

export const tabs: TabContent = {
  details: {
    name: "Details",
    icon: "documentDetails",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  yaml: {
    name: "YAML",
    icon: "code",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  new: {
    name: "New",
    icon: "plus",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  logs: {
    name: "Logs",
    icon: "documentList",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  shell: {
    name: "Shell",
    icon: "cmdLine",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  events: {
    name: "Events",
    icon: "newspaper",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  queue: {
    name: "Queue",
    icon: "queue",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  stack: {
    name: "Stack",
    icon: "stack",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  search: {
    name: "Search",
    icon: "search",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  dashboard: {
    name: "Dashboard",
    icon: "dashboard",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  repo: {
    name: "Repository",
    icon: "repo",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  list: {
    name: "List",
    icon: "list",
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
} as const;

export const routeString = {
  auth: "/auth",
  home: "/",
  helm: "/helm",
  nodeList: "/nodes",
  podList: "/pods",
  deployList: "/deployments",
  serviceList: "/services",
  secretList: "/secrets",
  configMapList: "/configmaps",
  statefulSetList: "/statefulsets",
  ingressList: "/ingresses",
  jobList: "/jobs",
  clusterRoleList: "/clusterroles",
  clusterRoleBindingList: "/clusterrolebindings",
  roleList: "/roles",
  roleBindingList: "/rolebindings",
  cronJobList: "/cronjobs",
  daemonSetList: "/daemonsets",
  replicaSetList: "/replicasets",
  eventList: "/events",
  ingressClassList: "/ingressclasses",
  namespaceList: "/namespaces",
  networkPolicyList: "/networkpolicies",
  persistentVolumeList: "/persistentvolumes",
  persistentVolumeClaimList: "/persistentvolumeclaims",
  replicationControllerList: "/replicationcontrollers",
  serviceAccountList: "/serviceaccounts",
  storageClassList: "/storageclasses",
} as const;

export const EmptyGVRK: KubeGVRK = {
  group: "",
  version: "",
  resource: "",
  kind: "",
  isNamespaced: false,
};

export const SelfSubjectAccessReviewV1GVRK: KubeGVRK = {
  group: "authorization.k8s.io",
  version: "v1",
  resource: "selfsubjectaccessreviews",
  kind: "SelfSubjectAccessReview",
  isNamespaced: false,
};

export const NamespaceV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "namespaces",
  kind: "Namespace",
  isNamespaced: false,
} as const;

export const NodeV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "nodes",
  kind: "Node",
  isNamespaced: false,
} as const;

export const NodeMetricsV1GVRK: KubeGVRK = {
  group: "metrics.k8s.io",
  version: "v1beta1",
  resource: "nodes",
  kind: "NodeMetrics",
  isNamespaced: false,
} as const;

export const PersistentVolumeV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "persistentvolumes",
  kind: "PersistentVolume",
  isNamespaced: false,
} as const;

export const PersistentVolumeClaimV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "persistentvolumeclaims",
  kind: "PersistentVolumeClaim",
  isNamespaced: true,
} as const;

export const ReplicationControllerV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "replicationcontrollers",
  kind: "ReplicationController",
  isNamespaced: true,
} as const;

export const PodV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "pods",
  kind: "Pod",
  isNamespaced: true,
} as const;

export const PodMetricsV1GVRK: KubeGVRK = {
  group: "metrics.k8s.io",
  version: "v1beta1",
  resource: "pods",
  kind: "PodMetrics",
  isNamespaced: true,
} as const;

export const ServiceV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "services",
  kind: "Service",
  isNamespaced: true,
} as const;

export const SecretV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "secrets",
  kind: "Secret",
  isNamespaced: true,
} as const;

export const ServiceAccountV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "serviceaccounts",
  kind: "ServiceAccount",
  isNamespaced: true,
} as const;

export const ConfigMapV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "configmaps",
  kind: "ConfigMap",
  isNamespaced: true,
} as const;

export const DeploymentV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "deployments",
  kind: "Deployment",
  isNamespaced: true,
} as const;

export const DaemonSetV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "daemonsets",
  kind: "DaemonSet",
  isNamespaced: true,
} as const;

export const ReplicaSetV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "replicasets",
  kind: "ReplicaSet",
  isNamespaced: true,
} as const;

export const StatefulSetV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "statefulsets",
  kind: "StatefulSet",
  isNamespaced: true,
} as const;

export const EventV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "events",
  kind: "Event",
  isNamespaced: true,
} as const;

export const JobV1GVRK: KubeGVRK = {
  group: "batch",
  version: "v1",
  resource: "jobs",
  kind: "Job",
  isNamespaced: true,
} as const;

export const CronJobV1GVRK: KubeGVRK = {
  group: "batch",
  version: "v1",
  resource: "cronjobs",
  kind: "CronJob",
  isNamespaced: true,
} as const;

export const IngressV1GVRK: KubeGVRK = {
  group: "networking.k8s.io",
  version: "v1",
  resource: "ingresses",
  kind: "Ingress",
  isNamespaced: true,
} as const;

export const IngressClassV1GVRK: KubeGVRK = {
  group: "networking.k8s.io",
  version: "v1",
  resource: "ingressclasses",
  kind: "Ingressclass",
  isNamespaced: false,
} as const;

export const NetworkPolicyV1GVRK: KubeGVRK = {
  group: "networking.k8s.io",
  version: "v1",
  resource: "networkpolicies",
  kind: "NetworkPolicy",
  isNamespaced: true,
} as const;

export const RoleV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "roles",
  kind: "Role",
  isNamespaced: true,
} as const;

export const RoleBindingV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "rolebindings",
  kind: "RoleBinding",
  isNamespaced: true,
} as const;

export const ClusterRoleV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "clusterroles",
  kind: "ClusterRole",
  isNamespaced: false,
} as const;

export const ClusterRoleBindingV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "clusterrolebindings",
  kind: "ClusterRoleBinding",
  isNamespaced: false,
} as const;

export const StorageClassV1GVRK: KubeGVRK = {
  group: "storage.k8s.io",
  version: "v1",
  resource: "storageclasses",
  kind: "StorageClass",
  isNamespaced: false,
} as const;

export const CustomResourceDefinitionV1GVRK: KubeGVRK = {
  group: "apiextensions.k8s.io",
  version: "v1",
  resource: "customresourcedefinitions",
  kind: "CustomResourceDefinition",
  isNamespaced: false,
} as const;

export const nanoToMiliCoreDivider = 1000000;

export const WSCloseCode = {
  info: 3000,
  warning: 3001,
  error: 3002,
  refresh: 3003,
} as const;

export const EventType = {
  NORMAL: "Normal",
  WARNING: "Warning",
} as const;

export const KubeAuthType = {
  inClusterConfig: "inClusterConfig",
  kubeconfigPath: "kubeconfigPath",
  kubeconfigRaw: "kubeconfigRaw",
  accessToken: "accessToken",
} as const;

export const transitionEffects: string =
  "transition-all duration-150 ease-in-out";

export const jsonStringClassName: string =
  "max-h-fit w-full rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm whitespace-pre";

export function getLabels(
  labels: { [key: string]: string } | undefined,
): EmbeddedTableItem[] {
  let arr: EmbeddedTableItem[] = [];
  if (labels !== undefined)
    Object.entries(labels).forEach((label) => {
      arr.push({ name: label[0], value: label[1] });
    });
  return arr;
}

export function getContainers(containers: V1Container[]): string[] {
  return containers.map((c) => {
    return c.name;
  });
}

export function parseFieldSelector(selector: KubeFieldSelector): string {
  return `involvedObject.kind=${selector.kind},involvedObject.name=${selector.name},involvedObject.namespace=${selector.namespace}`;
}

export function randomUUID(): string {
  return self.crypto.randomUUID();
}
