import { writable } from "svelte/store";
import type { SettingsData, SidebarState } from "./types";

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

export const availableContainers = writable<Array<string>>([]);

export const embeddedActiveTab = writable<number | undefined>();

export const yamlHasChanges = writable<boolean>(false);
export const discardYamlChanges = writable<boolean>(false);
export const submitYamlChanges = writable<boolean>(false);

export const statusCode = writable<number>();
export const kubeAuthState = writable<string>(
  localStorage.getItem("authState") ?? "false",
);
