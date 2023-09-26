<script lang="ts">
  import DialogElement from "../DialogElement.svelte";
  import SvgIcon from "../SvgIcon.svelte";
  import type { DialogData, EmbeddedOptionsData } from "../types";
  import { link } from "svelte-spa-router";

  export let embeddedOptionsData: EmbeddedOptionsData[];
  export let optClassNames: string = "";

  let el: Element;
  let showDialog = false;

  $: elCount = el?.children?.length;
  $: dialogData = {} as DialogData;
</script>

<DialogElement
  bind:showDialog
  dialogData={{
    type: dialogData.type,
    action: dialogData.action,
    resourceName: dialogData.resourceName,
  }}
/>

<div
  class="dropdown dropdown-left dropdown-hover hover:bg-base-200 flex
    h-full items-center space-x-2 rounded-lg p-1 {optClassNames}"
>
  <SvgIcon type={"elipsisVertical"} classNames={"size-4 h-full"} />

  <div
    bind:this={el}
    class="dropdown-content z-[1] flex h-full place-items-center items-center pr-2
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
            <SvgIcon type={option.icon} />
          </a>
        {:else}
          <button
            class="btn join-item btn-xs {option.classes}"
            on:click={() => {
              option.fn();
              if (option.dialog) {
                showDialog = true;
                dialogData = option.dialog;
              }
            }}
          >
            <SvgIcon type={option.icon} />
          </button>
        {/if}
      {/if}
    {/each}
  </div>
</div>
