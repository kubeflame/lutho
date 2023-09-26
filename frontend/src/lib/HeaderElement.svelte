<script lang="ts">
  import ClusterConnection from "./ClusterConnection.svelte";
  import NamespaceSelect from "./NamespaceSelect.svelte";
  import NavControl from "./NavControl.svelte";
  import SvgIcon from "./SvgIcon.svelte";
  import Utils from "./Utils.svelte";
  import { kubeHost, showUtils, sidebarState } from "./stores";
  import { transitionEffects } from "./util";

  export let headerContentClassName: string = "";
  export let logoClassName: string = "";
  export let showNamespaceSelect: boolean = false;

  const logoName = "LUTHO";
  const fontLogo = "font-mono font-thin tracking-wider";

  let connectionDialog: HTMLDialogElement;

  $: headerUtilsStyle = $kubeHost
    ? "border-base-300 bg-base-100 text-center"
    : "text-base-content/55 border-error bg-base-200 text-center";
</script>

<ClusterConnection bind:dialog={connectionDialog} />

<Utils />

<div
  class="header-content no-scrollbar {headerContentClassName} bg-base-200
    flex h-10 items-center justify-center p-1"
>
  <NavControl />

  <div
    class="logo {logoClassName} flex px-6 {transitionEffects} {fontLogo}
      {$sidebarState === 'max' ? '-ml-20' : 'ml-10'}"
  >
    {logoName}
  </div>

  <div
    class="mb-1 grid grid-flow-col grid-rows-1 items-center
      justify-center gap-6 {transitionEffects}"
  >
    <slot name="tabs" />
  </div>

  <div class="spacer flex grow items-center" />

  <div
    class="header-utils mr-2 flex !h-7 items-center space-x-2 {transitionEffects}"
  >
    {#if showNamespaceSelect}
      <NamespaceSelect />
    {/if}

    <button
      class="btn-xs bg-base-100 border-base-300 hover:bg-primary flex
        !h-full items-center rounded-xl border p-3 shadow-sm"
      on:click={() => ($showUtils = !$showUtils)}
    >
      <SvgIcon type={"util"} classNames={"size-5"} />
    </button>

    <div
      class="join flex !h-full content-center items-center
        rounded-xl font-mono text-sm font-thin tracking-wide shadow-sm"
    >
      <button
        class="btn join-item btn-xs bg-base-100 hover:bg-primary
          flex h-full min-w-fit max-w-fit {headerUtilsStyle}"
        on:click={() => connectionDialog.showModal()}
      >
        <SvgIcon
          type={$kubeHost ? "cloudOn" : "cloudOff"}
          fill={$kubeHost ? "oklch(var(--p))" : "oklch(var(--er))"}
          stroke={"oklch(var(--b1))"}
          classNames={"size-6"}
          strokeWidth={1.5}
        />
      </button>
      <div
        id="cluster-connection"
        class="dropdown dropdown-hover join-item relative flex h-full max-w-32
          content-center items-center border px-2 lg:max-w-md xl:max-w-lg {headerUtilsStyle}"
      >
        <span class="overflow-hidden text-ellipsis">
          {$kubeHost || "not connected"}
        </span>
      </div>
    </div>
  </div>
</div>
