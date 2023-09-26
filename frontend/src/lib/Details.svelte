<script lang="ts">
  import SvgIcon from "./SvgIcon.svelte";
  import { transitionEffects } from "./util";

  export let title: string;
  export let contentClassNames: string = "";
</script>

<div class="mx-1 mb-4 mt-2">
  <div class="collapse-plus collapse rounded-lg shadow-md {transitionEffects}">
    <input type="checkbox" class="peer h-8" checked />

    <div
      class="collapse-title text-base-content/80 bg-base-200 flex h-8 !items-center rounded-t-lg
        pl-3 font-mono text-sm font-light tracking-wide after:!right-4 after:!top-1.5"
    >
      <SvgIcon type={"info"} />
      <p class="ml-2 flex items-center">{title}</p>
    </div>

    <div
      class="collapse-content space-y-2 overflow-hidden rounded-b-md peer-checked:!pt-4
        {contentClassNames}"
    >
      <slot />
    </div>
  </div>
</div>
