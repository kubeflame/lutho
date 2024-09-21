<script lang="ts">
  import { slide } from "svelte/transition";
  import SvgIcon from "./SvgIcon.svelte";
  import type { Alert } from "./types";
  import type { icons } from "./icons";
  import { sleep } from "./util";
  import { onDestroy } from "svelte";

  export let alert: Alert;

  const alertClassNames = "z-20 rounded-lg absolute bottom-0";
  const svgIconStrokeWidth: number = 1.5;
  const svgIconClassNames: string = "h-6 w-6 stroke-current";

  const alertColor = (at: Alert["type"]): string => {
    switch (at) {
      case "info":
        return "alert-info";
      case "success":
        return "alert-success";
      case "warning":
        return "alert-warning";
      case "error":
        return "alert-error";
      default:
        return "";
    }
  };

  $: iconType = alert.type as keyof typeof icons;
  $: if (alert.timeout) {
    sleep(alert.timeout).then(() => {
      alert.message = null;
      alert.type = null;
      alert.timeout = null;
    });
  }

  onDestroy(() => {
    alert.message = null;
    alert.type = null;
    alert.timeout = null;
  });
</script>

{#if alert.message}
  <div
    class="alert {alertColor(alert.type)} {alert.className || alertClassNames}"
    in:slide={{ duration: 200 }}
    out:slide={{ duration: 200 }}
  >
    <SvgIcon classNames={svgIconClassNames} type={iconType} />
    <span>{alert.message}</span>
    <div class="spacer flex grow items-center" />
    <button
      class="btn btn-circle btn-ghost btn-sm"
      on:click={() => {
        alert.message = null;
        alert.type = null;
        alert.timeout = null;
      }}
    >
      <SvgIcon strokeWidth={svgIconStrokeWidth} type={"close"} />
    </button>
  </div>
{/if}
