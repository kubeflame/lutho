<script lang="ts">
  import { fade } from "svelte/transition";
  import SvgIcon from "../SvgIcon.svelte";
  import type { DashboardStat } from "../types";

  export let stats: DashboardStat[] = [];
</script>

<div class="stats mx-1 flex flex-1 rounded-lg shadow-md" transition:fade>
  {#each stats as stat}
    <div class="stat flex flex-col place-content-center items-center gap-y-1">
      <div class="stat-title flex items-center">
        <SvgIcon type={stat.iconType} />
        <span class="mx-1 flex">{stat.title}</span>
      </div>
      <div class="stat-value flex items-center">{stat.value}</div>
      <div class="stat-desc flex items-center">{stat.description}</div>
    </div>
  {/each}
</div>
