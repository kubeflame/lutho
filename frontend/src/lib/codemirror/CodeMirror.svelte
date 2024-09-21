<script lang="ts">
  import CodeMirror from "svelte-codemirror-editor";
  import { cmTheme, cmHighlightYAML } from "./theme";
  import { EditorView, ViewUpdate } from "@codemirror/view";
  import type { Alert } from "../types";
  import { load as fromYaml } from "js-yaml";
  import BottomAlert from "../BottomAlert.svelte";
  import { discardYamlChanges, settingsData } from "../stores";
  import { writable } from "svelte/store";
  import { undo } from "@codemirror/commands";
  import { transitionEffects } from "../util";
  import { basicSetup } from "./basicExtensions";

  export let value = "";
  export let codeMirrorChanged = false;
  export let docStore = writable<string>();

  let errorMessage: string = "";
  let view: EditorView;

  $: extensions = [
    basicSetup,
    cmTheme,
    EditorView.updateListener.of((v: ViewUpdate) => {
      codeMirrorChanged = v.state.doc.toString() !== value;
    }),
  ];
  $: alert = { message: errorMessage, type: "error" } as Alert;
  $: $settingsData.codemirrorExtraExtensions.yamlHighlight === true
    ? (extensions = [extensions, cmHighlightYAML])
    : removeExtension(cmHighlightYAML);

  $: if ($discardYamlChanges) {
    handleUndo();
    errorMessage = "";
  }

  function checkYamlForErrors(val: string) {
    try {
      fromYaml(val);
      errorMessage = "";
    } catch (error) {
      errorMessage = error as string;
    }
  }

  function removeExtension(extension: any) {
    extensions = extensions.filter((ex) => {
      return ex !== extension;
    });
  }

  function handleUndo() {
    if (view) {
      while (codeMirrorChanged)
        undo({
          state: view.state,
          dispatch: view.dispatch,
        });
      $discardYamlChanges = false;
    }
  }
</script>

<CodeMirror
  {value}
  {extensions}
  basic={false}
  lineWrapping
  nodebounce
  on:change={(e) => {
    $docStore = e.detail;
    checkYamlForErrors(e.detail);
  }}
  on:ready={(e) => {
    if ($docStore) {
      e.detail.dispatch({
        changes: {
          from: 0,
          to: e.detail.state.doc.length,
          insert: $docStore,
        },
      });
    }
    view = e.detail;
  }}
  class="codemirror flex h-full transform !pr-2 {transitionEffects}"
/>

<BottomAlert bind:alert />
