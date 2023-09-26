<script lang="ts">
  import { sidebarState } from "./stores";
  import { location, link } from "svelte-spa-router";
  import type { SidebarItem } from "./types";
  import SvgIcon from "./SvgIcon.svelte";
  import { routeString, transitionEffects } from "./util";

  const svgIconClassNames: string = "h-6 w-6";

  const sidebarItems: SidebarItem[] = [
    {
      name: "Home",
      icon: "home",
      url: "home",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Helm",
      icon: "helm",
      url: "helm",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Nodes",
      icon: "node",
      url: "nodeList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Namespaces",
      icon: "namespace",
      url: "namespaceList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Pods",
      icon: "pod",
      url: "podList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Ingresses",
      icon: "ingress",
      url: "ingressList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Deployments",
      icon: "deployment",
      url: "deployList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Services",
      icon: "service",
      url: "serviceList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Secrets",
      icon: "secret",
      url: "secretList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "ConfigMaps",
      icon: "configMap",
      url: "configMapList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "StatefulSets",
      icon: "queue",
      url: "statefulSetList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Jobs",
      icon: "job",
      url: "jobList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "ClusterRoles",
      icon: "clusterRole",
      url: "clusterRoleList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "ClusterRoleBindings",
      icon: "clusterRoleBinding",
      url: "clusterRoleBindingList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Roles",
      icon: "role",
      url: "roleList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "RoleBindings",
      icon: "roleBinding",
      url: "roleBindingList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "CronJobs",
      icon: "cronjob",
      url: "cronJobList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "DaemonSets",
      icon: "daemonset",
      url: "daemonSetList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "ReplicaSets",
      icon: "hashtag",
      url: "replicaSetList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "Events",
      icon: "newspaper",
      url: "eventList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "IngressClasses",
      icon: "hashtag",
      url: "ingressClassList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "NetworkPolicies",
      icon: "hashtag",
      url: "networkPolicyList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "PersistentVolumes",
      icon: "volume",
      url: "persistentVolumeList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "PersistentVolumeClaims",
      icon: "persistentVolumeClaim",
      url: "persistentVolumeClaimList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "ReplicationController",
      icon: "hashtag",
      url: "replicationControllerList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "ServiceAccounts",
      icon: "serviceAccount",
      url: "serviceAccountList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
    {
      name: "StorageClasses",
      icon: "storageClass",
      url: "storageClassList",
      viewBox: "0 0 24 24",
      strokeWidth: 1,
    },
  ];
</script>

{#if $location !== routeString.auth}
  <aside
    class="aside no-scrollbar bg-base-200 fixed bottom-0 left-0 top-10 z-50 h-screen
    {$sidebarState === 'min' ? 'w-12' : 'w-60'} {transitionEffects}"
  >
    <div
      class="max-sidebar no-scrollbar absolute flex h-[96vh] flex-col gap-y-1
      overflow-y-auto overflow-x-hidden py-2 {transitionEffects}"
    >
      {#each sidebarItems as item}
        <a
          role="button"
          href={routeString[item.url]}
          use:link
          class="btn btn-sm hover:bg-primary border-base-100 hover:border-base-100
          z-50 flex place-items-center items-center justify-start rounded-e-3xl
          border-2 p-0.5 px-2 font-medium tracking-wide hover:pl-6

          {$location.split('/')[1] === routeString[item.url].split('/')[1]
            ? 'bg-primary'
            : 'bg-base-200'} 
          
          {$sidebarState === 'max'
            ? 'w-60 hover:w-64'
            : 'tooltip tooltip-right tooltip-primary w-12 hover:mr-80 hover:w-16'} 
          {transitionEffects}"
          data-tip={item.name}
        >
          <SvgIcon
            classNames={svgIconClassNames}
            strokeWidth={item.strokeWidth}
            viewBox={item.viewBox}
            type={item.icon}
          />
          {#if $sidebarState === "max"}
            <div class={transitionEffects}>{item.name}</div>
          {/if}
        </a>
      {/each}
    </div>
  </aside>
{/if}
