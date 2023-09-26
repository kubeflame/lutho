<script lang="ts">
  import { icons } from "./util";
  import { showSettings, settingsData } from "./stores";

  let codemirrorYamlHighlight: boolean =
    $settingsData.codemirrorExtraExtensions?.yamlHighlight ?? true;

  $: {
    $settingsData = {
      codemirrorExtraExtensions: {
        yamlHighlight: codemirrorYamlHighlight,
      },
    };
    localStorage.setItem("settingsData", JSON.stringify($settingsData));
  }

  let settingsDialog: HTMLDialogElement;
  $: if (settingsDialog && $showSettings) settingsDialog.showModal();
</script>

<dialog bind:this={settingsDialog} class="modal modal-middle">
  <div
    class="h-full w-full rounded-2xl border-4 border-base-100 bg-base-200 drop-shadow-xl"
  >
    <div class="m-2">
      <button
        class="bg-none hover:bg-none hover:text-error"
        on:click={() => {
          $showSettings = false;
          settingsDialog.close();
        }}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
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

    <div class="place-items-cente m-2 ml-6 flex items-center">
      <div class="flex place-content-center items-center">
        Enable YAML highlighting <input
          type="checkbox"
          class="toggle toggle-success toggle-sm ml-2"
          bind:checked={codemirrorYamlHighlight}
        />
      </div>
    </div>
  </div>
</dialog>
