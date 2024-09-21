import { writable } from "svelte/store";
import type { SettingsData, SidebarState } from "./types";
import type { V1SelfSubjectAccessReview } from "@kubernetes/client-node";

export const namespace = writable<string>(
  localStorage.getItem("namespace") ?? "default",
);

export const sidebarState = writable<SidebarState>(
  (localStorage.getItem("navState") as SidebarState) ?? "max",
);

export const showSettings = writable<boolean>(false);
export const settingsData = writable<SettingsData>(
  JSON.parse(localStorage.getItem("settingsData") as string) ??
    ({} as SettingsData),
);

export const showUtils = writable<boolean>(false);

export const yamlHasChanges = writable<boolean>(false);
export const discardYamlChanges = writable<boolean>(false);
export const submitYamlChanges = writable<boolean>(false);

export const statusCode = writable<number>();

export const kubeHost = writable<string>("");
export const defaultAccessReview = writable<V1SelfSubjectAccessReview>();

export const updatingRepos = writable<{ [key: string]: boolean }>({});
