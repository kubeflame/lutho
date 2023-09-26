<script lang="ts">
  import type { Usage } from "@kubernetes/client-node";
  import CaptionTag from "./CaptionTag.svelte";
  import SvgIcon from "./SvgIcon.svelte";

  export let usageData: Usage = {} as Usage;
  export let strokeWidth: number = 1.5;
</script>

<div class="grid grid-cols-1 gap-5 xl:grid-cols-2">
  <div class="grid">
    <CaptionTag tagName={"Memory"} />

    <div class="stats flex-shrink shadow">
      <div class="stat">
        <div class="stat-figure">
          <SvgIcon
            type="decrease"
            classNames="size-6 text-primary"
            {strokeWidth}
          />
        </div>
        <div class="stat-value">{usageData.memory}</div>
      </div>
    </div>
  </div>

  <div class="grid">
    <CaptionTag tagName={"CPU"} />

    <div class="stats flex-shrink shadow">
      <div class="stat">
        <div class="stat-figure">
          <SvgIcon
            type="increase"
            classNames="size-6 text-error"
            {strokeWidth}
          />
          <SvgIcon
            type="steady"
            classNames="size-6 text-secondary"
            {strokeWidth}
          />
          <SvgIcon
            type="decrease"
            classNames="size-6 text-primary"
            {strokeWidth}
          />
        </div>
        <div class="stat-value">{usageData.cpu}</div>
      </div>
    </div>
  </div>
</div>
