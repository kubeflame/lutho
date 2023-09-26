<script lang="ts">
  import SvgIcon from "./SvgIcon.svelte";
  import { showSettings, settingsData } from "./stores";

  const svgIconClassNames: string = "h-6 w-6";
  const svgIconStrokeWidth: number = 1.5;

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

<dialog
  bind:this={settingsDialog}
  class="modal modal-middle"
  on:close={() => ($showSettings = false)}
>
  <div
    class="h-full w-full rounded-2xl border-4 border-base-100 bg-base-200 drop-shadow-xl"
  >
    <div class="m-2">
      <button
        class="bg-none hover:bg-none hover:text-error"
        on:click={() => settingsDialog.close()}
      >
        <SvgIcon
          classNames={svgIconClassNames}
          strokeWidth={svgIconStrokeWidth}
          type={"close"}
        />
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
