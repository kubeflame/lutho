export type HTTPStatusCode = {
  [key: number]: string;
};

export type ServerResponse = {
  sessionId: string;
  error: string;
  statusCode: number;
};

export type ServerRequest = {
  type?: any;
  url?: string;
  urlPArams?: string;
  group?: string;
  version?: string;
  resource?: string;
  kind?: string;
  signal?: AbortSignal;
  namespace?: string;
  name?: string;
  data?: any;
  options?: KubeOptions;
};

export type KubeOptions = {
  fieldSelector?: string;
  labelSelector?: string;
  timeoutSeconds?: number;
  limit?: number;
  _continue?: string;
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
    tags?: Array<string>;
    enabled?: boolean;
    "import-values"?: Array<any>;
    alias?: string;
  };

  export type Metadata = {
    name?: string;
    home?: string;
    sources: Array<string>;
    version?: string;
    description?: string;
    keywords: Array<string>;
    maintainers: Array<Maintainer>;
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
    dependencies: Array<Dependency>;
    type?: string;
  };

  export type Lock = {
    generated?: Date;
    digest?: string;
    dependencies: Array<Dependency>;
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
    events?: Array<string>;
    last_run?: HookExecution;
    wieght?: number;
    delete_policies?: Array<string>;
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
    "-": Array<File>;
    metadata?: Metadata;
    lock?: Lock;
    templates: Array<File>;
    values?: {
      [key: string]: any;
    };
    schema?: ArrayBuffer;
    files: Array<File>;
    parent?: Chart;
    depdendencies: Array<Chart>;
  };

  export type Release = {
    name?: string;
    info: ReleaseInfo;
    chart?: Chart;
    config?: {
      [key: string]: any;
    };
    manifest?: string;
    hooks?: Array<Hook>;
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

export type XtermData = {
  namespace: string;
  name: string;
  containers?: string;
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
  type: string;
  request?: ServerRequest;
};

export const KubeDataOpType = {
  bind: "bind",
  list: "list",
  helmList: "helmList",
  get: "get",
  helmGet: "helmGet",
  check: "check",
  update: "update",
  delete: "delete",
  close: "close",
  stdin: "stdin",
  stdout: "stdout",
  resize: "resize",
  toast: "toast",
} as const;

export type Alert = {
  type: "info" | "warning" | "error" | null;
  message: string | null;
};

export const WSCloseCode = {
  info: 3000,
  warning: 3001,
  error: 3002,
  refresh: 3003,
} as const;

export const TabIndex = {
  DETAILS: 0,
  YAML: 1,
  LOGS: 2,
  SHELL: 3,
  EVENTS: 4,
  QUEUE: 5,
  STACK: 6,
  SHELL_EMBED: 7,
  LOGS_EMBED: 8,
  SEARCH: 9,
  DASHBOARD: 10,
} as const;

export type TabItem = {
  index: number;
  name: string;
  icon: string;
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
  icon: string;
  url: string;
  viewBox: string;
  strokeWidth: number;
};

export const EventType = {
  NORMAL: "Normal",
  WARNING: "Warning",
} as const;

export type CodeMirrorExtraExtensions = {
  yamlHighlight: boolean;
};

export type SettingsData = {
  codemirrorExtraExtensions: CodeMirrorExtraExtensions;
};

export type EmbeddedOptionsData = {
  icon: string;
  classes: string;
  fn: Function;
  url?: string;
  hide?: boolean;
};

export type EmbeddedTableItem = {
  name: string;
  value: string | undefined;
  hidden?: boolean;
};

export type KubeGVRK = {
  group: string;
  version: string;
  resource: string;
  kind: string;
  isNamespaced?: boolean;
};
