<script lang="ts">
  import BottomAlert from "./BottomAlert.svelte";
  import SvgIcon from "./SvgIcon.svelte";
  import type { Alert, DialogData, DialogStep } from "./types";

  export let dialogData: DialogData;
  export let showDialog: boolean;
  export let alert: Alert = { message: null, type: null };
  export let onClose: () => void = () => {};
  export let isLoading: boolean = false;
  export let step: DialogStep = "action";

  let dialogElement: HTMLDialogElement;

  $: if (dialogElement && showDialog) dialogElement.showModal();
</script>

<dialog
  bind:this={dialogElement}
  class="modal modal-middle {dialogData.classNames ?? ''}"
  on:close={() => {
    onClose();
    showDialog = false;
  }}
>
  <div
    class="modal-box grid max-h-[calc(70vh)] max-w-[calc(50vw)]
      grid-cols-1 rounded-lg border border-base-200 bg-base-100
      p-2 outline outline-1 outline-base-200 drop-shadow-md"
  >
    {#if step === "action"}
      <div class="text-md mb-5 flex items-center gap-x-1">
        <SvgIcon
          classNames={"h-5 w-5"}
          strokeWidth={1.5}
          type={dialogData.icon}
        />
        Are you sure you want to {dialogData.type?.toLocaleLowerCase()} '{dialogData.resourceName}'?
      </div>

      <div class="modal-action flex gap-x-2">
        <button
          class="btn btn-xs outline outline-1 outline-base-200 drop-shadow-sm
            hover:bg-warning focus:outline-1 focus:outline-warning"
          on:click={dialogData.action}
        >
          {dialogData.type}
        </button>

        <button
          class="btn btn-xs outline outline-1 outline-base-200 drop-shadow-sm
            hover:bg-error focus:outline-1 focus:outline-error"
          on:click={() => dialogElement.close()}
        >
          Close
        </button>
      </div>
    {:else if step === "afterAction"}
      <slot name="dialogContent" />

      <div class="modal-action flex gap-x-2">
        <slot name="afterActionLinks" />

        <div class="spacer flex grow items-center" />

        <button
          class="btn btn-xs outline outline-1 outline-base-200 drop-shadow-sm
          hover:bg-error focus:outline-1 focus:outline-error"
          on:click={() => dialogElement.close()}
        >
          Close
        </button>
      </div>
    {/if}

    <BottomAlert bind:alert />

    {#if isLoading}
      <progress
        class="progress progress-success mt-4 w-full ![animation-duration:_2s]"
      />
    {/if}
  </div>
</dialog>
