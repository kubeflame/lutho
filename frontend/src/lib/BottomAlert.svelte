<script lang="ts">
  import { slide } from "svelte/transition";
  import { icons } from "./util";
  import type { Alert } from "./types";

  export let alert: Alert;

  const iconType = (at: Alert["type"]): string | undefined => {
    switch (at) {
      case "info":
        return icons.info;
      case "warning":
        return icons.warning;
      case "error":
        return icons.closeCircle;
    }
  };

  const alertColor = (at: Alert["type"]): string | undefined => {
    switch (at) {
      case "info":
        return "alert-info";
      case "warning":
        return "alert-warning";
      case "error":
        return "alert-error";
    }
  };
</script>

{#if alert.message}
  <div
    class="alert absolute bottom-0 {alertColor(alert.type)} z-20 rounded-lg"
    in:slide={{ duration: 200 }}
    out:slide={{ duration: 200 }}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-6 w-6 shrink-0 stroke-current"
      fill="none"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="1.5"
        d={iconType(alert.type)}
      />
    </svg>
    <span>{alert.message}</span>
    <div class="spacer flex grow items-center" />
    <button
      class="btn btn-circle btn-ghost btn-sm"
      on:click={() => {
        alert.message = "";
      }}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-6 w-6 shrink-0 stroke-current"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d={icons.close}
        />
      </svg>
    </button>
  </div>
{/if}
