<script lang="ts">
  export let dialog: HTMLDialogElement;
</script>

<dialog bind:this={dialog} class="modal modal-middle z-50">
  <form
    method="dialog"
    class="modal-box bg-base-200 w-[calc(50vw)] max-w-[calc(50vw)]"
  >
    <h3 class="text-md mb-1 font-normal">Cluster connection</h3>
    <div class="modal-action">
      <div class="spacer flex grow items-center" />
      <button
        class="btn btn-xs outline-base-100 hover:bg-error focus:outline-error outline outline-1 drop-shadow-md focus:outline-1"
        on:click={() => dialog.close()}
      >
        Close
      </button>
    </div>
  </form>
</dialog>
