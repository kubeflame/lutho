<script lang="ts">
  import SvgIcon from "../SvgIcon.svelte";
  import type { DashboardStat } from "../types";
  import { transitionEffects } from "../util";

  export let stat: DashboardStat = {
    title: "empty",
    iconType: "hashtag",
    value: 0,
    description: "",
  };
  export let isLoading: boolean = false;
</script>

<div
  class="stats border-base-200 mx-1 flex flex-1 rounded-lg border shadow-sm {transitionEffects}"
>
  <div class="stat flex flex-col place-content-center items-center gap-y-1">
    <div class="stat-title flex items-center">
      <SvgIcon type={stat.iconType} />
      <span class="mx-1 flex">{stat.title}</span>
    </div>
    <div class="stat-value flex !h-10 items-center">
      {#if isLoading}
        <span class="loading loading-ring flex" />
      {:else}
        {stat.value}
      {/if}
    </div>
    <div class="stat-desc flex items-center">{stat.description}</div>
  </div>
</div>
