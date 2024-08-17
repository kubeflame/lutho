import type {
  V1CronJobList,
  V1DeploymentList,
  V1JobList,
  V1ListMeta,
  V1NodeList,
  V1PodList,
  V1ReplicaSetList,
  V1ReplicationControllerList,
  V1SelfSubjectAccessReview,
  V1ServiceList,
  V1StatefulSetList,
} from "@kubernetes/client-node";
import type { icons } from "./icons";
import type { KubeAuthType, routeString } from "./util";

export type HTTPStatusCode = {
  [key: number]: string;
};

export type ServerResponse = {
  sessionId: string;
  error: string;
  statusCode: number;
};

export type AuthResponse = {
  error: string;
  state: boolean;
  kubeHost: string;
  accessReview: V1SelfSubjectAccessReview;
};

export type AuthRequest = {
  type: keyof typeof KubeAuthType;
  token?: string;
  kubeconfigPath?: string;
  kubeconfigRaw?: string;
  masterURL?: string;
};

export type ServerRequest = {
  type?: any;
  url?: string;
  urlParams?: string;
  kubeGVRK?: KubeGVRK;
  kubeGVRKList?: KubeGVRK[];
  signal?: AbortSignal;
  namespace?: string;
  name?: string;
  data?: any;
  kubeOptions?: KubeOptions;
  helmOptions?: HelmOptions;
};

export type SockState = {
  state: number;
  bound: boolean;
  refresh?: boolean;
  query?: string;
};

export type KubeOptions = {
  fieldSelector?: string;
  labelSelector?: string;
  timeoutSeconds?: number;
  limit?: number;
  continue?: string;
};

export type HelmOptions = {
  chartName?: string;
  chartVersion?: string;
  envPath?: string;
  repoURL?: string;
  dryRun?: boolean;
  isOCI?: boolean;
  reuseValues?: boolean;
};

export namespace KubeEvent {
  export type Metadata = {
    creationTimestamp?: Date;
    name: string;
    namespace: string;
  };

  export type Regarding = {
    apiVersion: string;
    kind: string;
    name: string;
    namespace: string;
  };

  export type Series = {
    count?: number;
    lastObservedTime?: Date;
  };

  export type Event = {
    action: string;
    kind: string;
    metadata: Metadata;
    regarding: Regarding;
    note: string;
    type: string;
    series: Series;
    reason: string;
    eventTime: Date;
    reportingController?: string;
    reportingInstance?: string;
  };
}

export namespace Helm {
  export type Maintainer = {
    name?: string;
    email?: string;
    url?: string;
  };

  export type Dependency = {
    name?: string;
    version?: string;
    repository?: string;
    condition?: string;
    tags?: string[];
    enabled?: boolean;
    "import-values"?: any[];
    alias?: string;
  };

  export type Metadata = {
    name?: string;
    home?: string;
    sources: string[];
    version?: string;
    description?: string;
    keywords: string[];
    maintainers: Maintainer[];
    icon?: string;
    apiVersion?: string;
    condition?: string;
    tags?: string;
    appVersion?: string;
    deprecated?: boolean;
    annotations: {
      [key: string]: string;
    };
    kubeVersion?: string;
    dependencies: Dependency[];
    type?: string;
  };

  export type Lock = {
    generated?: Date;
    digest?: string;
    dependencies: Dependency[];
  };

  export type File = {
    name?: string;
    data?: string;
  };

  export type HookExecution = {
    started_at?: Date;
    completed_at?: Date;
    phase: string;
  };

  export type Hook = {
    name?: string;
    kind?: string;
    path?: string;
    manifest?: string;
    events?: string[];
    last_run?: HookExecution;
    wieght?: number;
    delete_policies?: string[];
  };

  export type ReleaseInfo = {
    first_deployed?: Date;
    last_deployed?: Date;
    deleted?: string;
    description?: string;
    status?: string;
    notes?: string;
    resources?: {
      [key: string]: any;
    };
  };

  export type Chart = {
    "-": File[];
    metadata?: Metadata;
    lock?: Lock;
    templates: File[];
    values?: {
      [key: string]: any;
    };
    schema?: ArrayBuffer;
    files: File[];
    parent?: Chart;
    depdendencies: Chart[];
  };

  export type Release = {
    name?: string;
    info: ReleaseInfo;
    chart?: Chart;
    config?: {
      [key: string]: any;
    };
    manifest?: string;
    hooks?: Hook[];
    version?: number;
    namespace?: string;
    "-"?: {
      [key: string]: string;
    };
  };
}

export type TerminalMessage = {
  op: string;
  data: string;
  sessionId: string;
  rows: number;
  cols: number;
  error?: string;
  statusCode?: number;
};

export type LogStreamMessage = {
  op: string;
  data: string;
  sessionId: string;
  error?: string;
  statusCode?: number;
};

export type KubeDataStreamMessage = {
  op: KubeOp;
  data?: any;
  sessionId: string;
  error?: string;
  statusCode?: number;
};

export type KubeOp = {
  type:
    | "bind"
    | "selfSubjectAccessReview"
    | "list"
    | "listAll"
    | "helmShowValues"
    | "helmList"
    | "get"
    | "helmGet"
    | "helmInstall"
    | "helmUpgrade"
    | "helmPull"
    | "helmGetTags"
    | "check"
    | "update"
    | "delete"
    | "helmUninstall"
    | "close"
    | "stdin"
    | "stdout"
    | "resize"
    | "toast";
  request?: ServerRequest;
  opID?: string;
};

export type KubeGVRK = {
  group: string;
  version: string;
  resource: string;
  kind: string;
  isNamespaced?: boolean;
};

export type Alert = {
  type: "info" | "warning" | "error" | null;
  message: string | null;
  className?: string;
};

export type TabItem = {
  name: string;
  icon: keyof typeof icons;
  strokeWidth: number;
  viewBox: string;
};

export class TabContent {
  [key: string]: TabItem;
}

export type ResourceToolbarContent = {
  index: number;
  name: string;
  url?: string;
};

export type SidebarState = "max" | "min";

export type SidebarItem = {
  name: string;
  icon: keyof typeof icons;
  url: keyof typeof routeString;
  viewBox: string;
  strokeWidth: number;
};

export type CodeMirrorExtraExtensions = {
  yamlHighlight: boolean;
};

export type SettingsData = {
  codemirrorExtraExtensions: CodeMirrorExtraExtensions;
};

export type HelmRepoData = {
  name: string;
  url: string;
  isOCI: boolean;
  indexFile?: string;
  latestVersion?: string;
  allVersions?: string[];
};

export type HelmChartTags = {
  chartName: string;
  chartTags: string[];
};

export type HelmChartIndexEntry = {
  name: string;
  version: string;
  created: string;
  dependencies: HelmChartIndexEntryDependency[];
  appVersion?: string;
  description?: string;
  kubeVersion?: string;
  maintainers?: any;
  sources?: any;
  urls?: any;
};

export type HelmChartIndexEntryDependency = {
  alias: string;
  condition: string;
  name: string;
  repository: string;
  version: string;
};

export type EmbeddedOptionsData = {
  icon: keyof typeof icons;
  classes: string;
  fn: (data?: any) => void;
  dialog?: DialogData;
  url?: string;
  hide?: boolean;
};

export type DialogData = {
  type?: string;
  action?: any;
  classNames?: string;
  resourceName?: string;
  icon?: keyof typeof icons;
  message?: string;
};

export type EmbeddedTableItem = {
  name: string;
  value: any;
  hidden?: boolean;
  className?: string;
};

export type TableType =
  | "normal"
  | "events"
  | "badges"
  | "one_row"
  | "custom-body"
  | "custom-vertical";

export type TabQueryParam =
  | "referrer"
  | "details"
  | "yaml"
  | "new"
  | "logs"
  | "shell"
  | "events"
  | "queue"
  | "stack"
  | "search"
  | "dashboard"
  | "repository"
  | "list";

export interface V1KubeListMeta extends V1ListMeta {
  continue: string;
}

export type KubeFieldSelector = {
  kind: string;
  name: string;
  namespace: string;
};

export type DashboardList = {
  [key: string]:
    | V1NodeList
    | V1PodList
    | V1DeploymentList
    | V1StatefulSetList
    | V1JobList
    | V1CronJobList
    | V1ServiceList
    | V1ReplicaSetList
    | V1ReplicationControllerList;
};

export type DashboardStat = {
  title: string;
  iconType: keyof typeof icons;
  value: number | null;
  description: string;
};
