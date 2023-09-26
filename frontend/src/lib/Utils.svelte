<script lang="ts">
  import BottomAlert from "./BottomAlert.svelte";
  import { showUtils } from "./stores";
  import type { Alert } from "./types";

  let dialog: HTMLDialogElement;
  let alert = {
    message: "",
    type: "error",
    className: "flex mt-2",
  } as Alert;

  $: value = "";
  $: if (dialog && $showUtils) dialog.showModal();

  function encode() {
    try {
      value = btoa(value);
      alert.message = "";
    } catch (error: any) {
      alert.message = error.message;
    }
  }

  function decode() {
    try {
      value = atob(value);
      alert.message = "";
    } catch (error: any) {
      alert.message = error.message;
    }
  }

  function clear() {
    value = "";
    alert.message = "";
  }
</script>

<dialog
  bind:this={dialog}
  class="modal modal-middle z-50"
  on:close={() => ($showUtils = false)}
>
  <div class="flex flex-col">
    <form
      method="dialog"
      class="modal-box bg-base-200 w-[calc(50vw)] max-w-[calc(50vw)]"
    >
      <h3 class="text-md mb-1 font-normal">
        Input the text to encode or decode using the <b>Base64</b> scheme:
      </h3>
      <textarea
        bind:value
        class="bg-base-200 outline-base-100 focus:outline-primary h-36 min-w-full resize-none overflow-y-scroll rounded-lg p-2 shadow-md outline outline-1"
      />
      <div class="modal-action">
        <button
          class="encode btn btn-xs outline-base-100 hover:bg-primary focus:outline-primary outline outline-1 drop-shadow-md focus:outline-1"
          on:click|preventDefault={encode}
        >
          Encode
        </button>
        <button
          class="decode btn btn-xs outline-base-100 hover:bg-primary focus:outline-primary outline outline-1 drop-shadow-md focus:outline-1"
          on:click|preventDefault={decode}
        >
          Decode
        </button>
        <div class="divider divider-horizontal" />
        <button
          class="decode btn btn-xs outline-base-100 hover:bg-primary focus:outline-primary outline outline-1 drop-shadow-md focus:outline-1"
          on:click|preventDefault={clear}
        >
          Clear
        </button>
        <button
          class="decode btn btn-xs outline-base-100 hover:bg-primary focus:outline-primary outline outline-1 drop-shadow-md focus:outline-1"
          on:click|preventDefault={() => {
            navigator.clipboard.writeText(value);
          }}
        >
          Copy
        </button>
        <button
          class="decode btn btn-xs outline-base-100 hover:bg-primary focus:outline-primary outline outline-1 drop-shadow-md focus:outline-1"
          on:click|preventDefault={() => {
            navigator.clipboard.readText().then((t) => (value = t));
          }}
        >
          Paste
        </button>
        <div class="spacer flex grow items-center" />
        <button
          class="btn btn-xs outline-base-100 hover:bg-error focus:outline-error outline outline-1 drop-shadow-md focus:outline-1"
          on:click={() => dialog.close()}
        >
          Close
        </button>
      </div>
    </form>
    <BottomAlert bind:alert />
  </div>
</dialog>
