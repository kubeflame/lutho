<script lang="ts">
  import NavControl from "./NavControl.svelte";
  import { kubeAuthState, sidebarState } from "./stores";
  import { icons } from "./util";

  export let headerContentClassName: string = "";
  export let logoClassName: string = "";

  const logoName = "LUTHO";
  const fontLogo = "font-mono font-thin tracking-wider";
  const userName = "Lutho Doe";

  function onClickAuthButton() {
    $kubeAuthState = $kubeAuthState === "false" ? "true" : "false";
    console.log("authState:", $kubeAuthState);
    localStorage.setItem("authState", $kubeAuthState);
  }
</script>

<header class="header sticky top-0 z-50">
  <div
    class="header-content {headerContentClassName} flex h-10 items-center
      justify-center bg-base-200 p-1"
  >
    <NavControl />

    <div
      class="logo {logoClassName} mt-1 flex px-5 duration-300 ease-in-out {fontLogo}
        {$sidebarState === 'max' ? '-ml-20 pr-8' : ''}"
    >
      {logoName}
    </div>

    <div
      class="grid transform grid-flow-col grid-rows-1 items-center justify-center
        gap-5 align-middle duration-300 ease-in-out"
    >
      <slot name="namespace" />
      <slot name="tabs" />
    </div>

    <div class="spacer flex grow items-center" />

    <div class="flex flex-none items-center justify-center text-center">
      <div class="flex items-center space-x-3 px-3">
        <div class="flex flex-none justify-center">
          <div class="flex h-8 w-8 items-center justify-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1"
              stroke="currentColor"
              class="rounded-full object-cover"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d={icons.user}
              />
            </svg>
          </div>
        </div>
        <div class="md:text-md hidden text-sm md:block">{userName}</div>
      </div>
    </div>
    <button class="btn btn-xs" on:click={onClickAuthButton}>Auth</button>
  </div>
  <slot name="toolbar" />
</header>
