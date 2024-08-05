<script lang="ts">
  import SvgIcon from "./SvgIcon.svelte";
  import type { DialogData } from "./types";

  export let dialogData: DialogData;
  export let showDialog: boolean;

  let dialogElement: HTMLDialogElement;

  $: if (dialogElement && showDialog) dialogElement.showModal();
</script>

<dialog
  bind:this={dialogElement}
  class="modal modal-middle {dialogData.classNames}"
  on:close={() => (showDialog = false)}
>
  <div
    class="modal-box border-base-100 bg-base-200 outline-base-200 grid max-h-[calc(70vh)]
      max-w-[calc(50vw)] grid-cols-1 rounded-lg border p-2 outline outline-1 drop-shadow-lg"
  >
    <div class="text-md mb-5 flex items-center gap-x-1">
      <SvgIcon
        classNames={"h-5 w-5"}
        strokeWidth={1.5}
        type={dialogData.icon}
      />
      Are you sure you want to {dialogData.type?.toLocaleLowerCase()} '{dialogData.resourceName}'?
    </div>
    <div class="modal-action">
      <form method="dialog">
        <button
          class="btn btn-xs outline-base-100 hover:bg-error focus:outline-error outline
            outline-1 drop-shadow-md focus:outline-1"
          on:click={dialogData.action}
        >
          {dialogData.type}
        </button>
        <button
          class="btn btn-xs outline-base-100 hover:bg-error focus:outline-error outline outline-1 drop-shadow-md focus:outline-1"
          on:click={() => dialogElement.close()}
        >
          Close
        </button>
      </form>
    </div>
  </div>
</dialog>
