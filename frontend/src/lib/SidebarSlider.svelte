<script lang="ts">
  import { sidebarState } from "./stores";
  import { location, link } from "svelte-spa-router";
  import type { SidebarItem } from "./types";
  import { slide } from "svelte/transition";
  import { cubicInOut } from "svelte/easing";

  export let sidebarItems: SidebarItem[];
</script>

{#if $sidebarState === "max"}
  <aside
    class="aside fixed bottom-0 left-0 top-10 z-10 h-screen w-60 bg-base-200 pl-0.5 pr-0.5"
    transition:slide={{ duration: 300, axis: "x", easing: cubicInOut }}
  >
    <div
      class="max-sidebar absolute top-2 flex h-[calc(100vh-60px)] flex-col gap-y-1
        overflow-y-auto overflow-x-hidden"
    >
      {#each sidebarItems as item}
        <a href={item.url} use:link class="max-link">
          <div
            class="max-item relative -left-4 z-10 flex w-64 transform place-items-center
            items-center gap-x-2 rounded-e-3xl border-2 border-base-100 bg-base-200 p-0.5
            pl-7 duration-200 ease-in-out hover:ml-3 hover:bg-primary hover:drop-shadow-sm
            {$location.split('/')[1] === item.url.split('/')[1]
              ? 'bg-primary'
              : 'bg-base-200'}"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox={item.viewBox}
              stroke-width={item.strokeWidth}
              stroke="currentColor"
              class="h-5 w-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d={item.icon}
              />
            </svg>
            <div>{item.name}</div>
          </div>
        </a>
      {/each}
    </div>
  </aside>
{:else if $sidebarState === "min"}
  <aside
    transition:slide={{ duration: 300, axis: "x", easing: cubicInOut }}
    class="aside fixed bottom-0 left-0 top-10 z-10 h-screen w-12 bg-base-200 pl-0.5 pr-0.5"
  >
    <div
      class="mini-sidebar absolute top-2 flex h-[calc(100vh-60px)] flex-col
        gap-y-1 overflow-y-auto overflow-x-hidden"
    >
      {#each sidebarItems as item}
        <a href={item.url} use:link class="mini-link">
          <div
            class="mini-item text-md tooltip-primary tooltip tooltip-right relative -left-4 z-10
              flex w-16 transform flex-row place-items-center items-center justify-end rounded-e-3xl
              border-2 border-base-100 bg-base-200 p-1 pr-3 font-medium tracking-wide duration-200
              ease-in-out hover:ml-3 hover:mr-96 hover:bg-primary hover:drop-shadow-sm
              {$location.split('/')[1] === item.url.split('/')[1]
              ? 'bg-primary'
              : 'bg-base-200'}"
            data-tip={item.name}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox={item.viewBox}
              stroke-width={item.strokeWidth}
              stroke="currentColor"
              class="h-5 w-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d={item.icon}
              />
            </svg>
          </div>
        </a>
      {/each}
    </div>
  </aside>
{/if}
