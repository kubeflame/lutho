<script lang="ts">
  import { scale } from "svelte/transition";
  import type { ResourceToolbarContent } from "./types";
  import { link } from "svelte-spa-router";
  import SvgIcon from "./SvgIcon.svelte";

  export let toolbarContent: ResourceToolbarContent[];
</script>

<div class="no-scrollbar breadcrumbs mx-4 text-sm">
  <ul>
    <div class="flex pl-0.5 pr-1.5">
      <SvgIcon type={"at"} />
    </div>
    {#each toolbarContent as item}
      <li in:scale={{ duration: 200 }}>
        {#if item.url}
          <a
            class="hover:text-primary drop-shadow-2xl hover:!no-underline"
            href={item.url}
            use:link
            on:click={() => {
              toolbarContent = toolbarContent.filter((x) => {
                return x.index <= item.index;
              });
              item.url = "";
            }}
          >
            {item.name}
          </a>
        {:else}
          {item.name}
        {/if}
      </li>
    {/each}
  </ul>
</div>

<div class="spacer flex grow items-center" />
