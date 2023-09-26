<script lang="ts">
  import { fade } from "svelte/transition";
  import { showUtils } from "./stores";

  let dialog: HTMLDialogElement;

  // const base64regex: RegExp =
  //     /^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$/;

  $: value = "";
  // $: b64 = value.match(base64regex);
  $: if (dialog && $showUtils) dialog.showModal();
</script>

<dialog bind:this={dialog} class="modal modal-middle z-50">
  <form
    method="dialog"
    class="modal-box w-[calc(50vw)] max-w-[calc(50vw)] bg-base-200"
  >
    <h3 class="text-md mb-1 font-normal">
      Input the text to encode or decode using the <b>Base64</b> scheme:
    </h3>
    <textarea
      bind:value
      class="h-36 min-w-full resize-none overflow-y-scroll rounded-lg bg-base-200 p-2 shadow-md outline outline-1 outline-base-100 focus:outline-primary"
    />
    <div class="modal-action">
      <button
        class="encode btn btn-xs outline outline-1 outline-base-100 drop-shadow-md hover:bg-primary focus:outline-1 focus:outline-primary"
        on:click|preventDefault={() => {
          try {
            value = btoa(value);
          } catch (error) {}
        }}
      >
        Encode
      </button>
      <button
        class="decode btn btn-xs outline outline-1 outline-base-100 drop-shadow-md hover:bg-primary focus:outline-1 focus:outline-primary"
        on:click|preventDefault={() => {
          try {
            value = atob(value);
          } catch (error) {}
        }}
      >
        Decode
      </button>
      <div class="divider divider-horizontal" />
      <button
        class="decode btn btn-xs outline outline-1 outline-base-100 drop-shadow-md hover:bg-primary focus:outline-1 focus:outline-primary"
        on:click|preventDefault={() => {
          value = "";
        }}
      >
        Clear
      </button>
      <button
        class="decode btn btn-xs outline outline-1 outline-base-100 drop-shadow-md hover:bg-primary focus:outline-1 focus:outline-primary"
        on:click|preventDefault={() => {
          navigator.clipboard.writeText(value);
        }}
      >
        Copy
      </button>
      <button
        class="decode btn btn-xs outline outline-1 outline-base-100 drop-shadow-md hover:bg-primary focus:outline-1 focus:outline-primary"
        on:click|preventDefault={() => {
          navigator.clipboard.readText().then((t) => (value = t));
        }}
      >
        Paste
      </button>
      <div class="spacer flex grow items-center" />
      <button
        class="btn btn-xs outline outline-1 outline-base-100 drop-shadow-md hover:bg-error focus:outline-1 focus:outline-error"
        on:click={() => {
          $showUtils = false;
        }}
      >
        Close
      </button>
    </div>
  </form>
</dialog>
