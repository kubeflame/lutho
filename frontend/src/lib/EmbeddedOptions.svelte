<script lang="ts">
  import { icons } from "./util";
  import type { EmbeddedOptionsData } from "./types";
  import { link } from "svelte-spa-router";

  export let embeddedOptionsData: EmbeddedOptionsData[];
  let el: Element;

  $: elCount = el?.children?.length;
</script>

<div
  class="dropdown-left dropdown-hover dropdown flex place-items-center items-center rounded-lg p-0.5 pl-0.5 hover:bg-base-200"
>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    stroke-width={1}
    stroke="currentColor"
    class="h-5 w-5"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      d={icons.elipsisVertical}
    />
  </svg>

  <div
    bind:this={el}
    class="dropdown-content z-[1] flex place-items-center items-center pr-2
        {elCount > 1 ? 'join' : ''}"
  >
    {#each embeddedOptionsData as option}
      {#if !option.hide}
        {#if option.url}
          <a
            use:link
            href={option.url}
            role="button"
            class="btn join-item btn-xs {option.classes}"
            on:click={() => {
              option.fn();
            }}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width={1}
              stroke="currentColor"
              class="h-4 w-4"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d={option.icon}
              />
            </svg>
          </a>
        {:else}
          <button
            class="btn join-item btn-xs {option.classes}"
            on:click={() => {
              option.fn();
            }}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width={1}
              stroke="currentColor"
              class="h-4 w-4"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d={option.icon}
              />
            </svg>
          </button>
        {/if}
      {/if}
    {/each}
  </div>
</div>
