<script lang="ts">
  import { slide } from "svelte/transition";
  import SvgIcon from "./SvgIcon.svelte";
  import type { Alert } from "./types";

  export let alert: Alert;

  const alertClassNames = "z-20 rounded-lg absolute bottom-0";
  const svgIconStrokeWidth: number = 1.5;
  const svgIconClassNames: string = "h-6 w-6 shrink-0 stroke-current";

  const iconType = (at: Alert["type"]) => {
    switch (at) {
      case "info":
        return "info";
      case "warning":
        return "warning";
      case "error":
        return "closeCircle";
    }
  };

  const alertColor = (at: Alert["type"]): string => {
    switch (at) {
      case "info":
        return "alert-info";
      case "warning":
        return "alert-warning";
      case "error":
        return "alert-error";
      default:
        return "";
    }
  };
</script>

{#if alert.message}
  <div
    class="alert {alertColor(alert.type)} {alert.className || alertClassNames}"
    in:slide={{ duration: 200 }}
    out:slide={{ duration: 200 }}
  >
    <SvgIcon classNames={svgIconClassNames} type={iconType(alert.type)} />
    <span>{alert.message}</span>
    <div class="spacer flex grow items-center" />
    <button
      class="btn btn-circle btn-ghost btn-sm"
      on:click={() => {
        alert.message = null;
        alert.type = null;
      }}
    >
      <SvgIcon strokeWidth={svgIconStrokeWidth} type={"close"} />
    </button>
  </div>
{/if}
