<script lang="ts">
  import { onMount } from "svelte";
  import {
    scrollPastEnd,
    lineNumbers,
    highlightActiveLine,
    highlightActiveLineGutter,
    EditorView,
    ViewUpdate,
  } from "@codemirror/view";
  import { EditorState, StateEffect } from "@codemirror/state";
  import { minimalSetup } from "codemirror";
  import {
    StreamLanguage,
    LanguageSupport,
    syntaxHighlighting,
  } from "@codemirror/language";
  import { classHighlighter } from "@lezer/highlight";
  import { yaml } from "@codemirror/legacy-modes/mode/yaml";
  import { discardYamlChanges, settingsData } from "./stores";
  import { writable } from "svelte/store";
  import BottomAlert from "./BottomAlert.svelte";
  import { load as fromYaml } from "js-yaml";
  import type { Alert } from "./types";
  import Utils from "./Utils.svelte";
  import { undo } from "@codemirror/commands";

  export let codeMirrorChanged = false;
  export let doc: string = "";
  export let docStore = writable<string>();

  let errorMessage: string = "";

  const yamlEditor = new LanguageSupport(StreamLanguage.define(yaml));
  let extensions = [
    minimalSetup,
    syntaxHighlighting(classHighlighter),
    scrollPastEnd(),
    lineNumbers(),
    EditorView.lineWrapping,
    highlightActiveLine(),
    highlightActiveLineGutter(),
    EditorView.updateListener.of((v: ViewUpdate) => {
      codeMirrorChanged = v.state.doc.toString() !== doc.toString();
      $docStore = v.state.doc.toString();
    }),
  ];

  $: try {
    fromYaml($docStore);
    errorMessage = "";
  } catch (error) {
    errorMessage = error as string;
  }

  $: alert = { message: errorMessage, type: "error" } as Alert;

  const handleUndo = () => {
    if (view) {
      while (codeMirrorChanged)
        undo({
          state: view.state,
          dispatch: view.dispatch,
        });
      $discardYamlChanges = false;
    }
  };

  $: if ($discardYamlChanges) {
    handleUndo();
    try {
      fromYaml($docStore);
      errorMessage = "";
    } catch (error) {
      errorMessage = error as string;
    }
  }

  let parent: HTMLDivElement;
  let view: EditorView;

  function _reconfigureExtensions() {
    if (view === null) return;
    view?.dispatch({
      effects: StateEffect.reconfigure.of(extensions),
    });
  }

  const removeExtension = (extension: any) => {
    extensions = extensions.filter((val) => {
      return val !== extension;
    });
  };

  $: extensions, _reconfigureExtensions();
  $: $settingsData.codemirrorExtraExtensions.yamlHighlight === true
    ? (extensions = [extensions, yamlEditor])
    : removeExtension(yamlEditor);

  onMount(() => {
    view = new EditorView({
      state: EditorState.create({
        doc,
        extensions: extensions,
      }),
      parent,
    });
  });
</script>

<Utils />

<div
  class="codemirror flex h-full transform bg-base-100 duration-300 ease-in-out"
  bind:this={parent}
>
  <BottomAlert bind:alert />
</div>

<style>
  :root {
    /* Dark color variables */
    --dark-tok-text: #cbd5e1;
    --dark-tok-comment: #64748b;
    --dark-tok-string: #fbbf24;
    --dark-tok-keyword: #a78bfa;
    --dark-tok-meta: #fb7185;
    --dark-tok-number: #10b981;
    --dark-tok-variable: #fb7185;
    --dark-tok-atom: #38bdf8;
    --dark-tok-def: #fb7185;

    /* Light color variables */
    --light-tok-text: #475569;
    --light-tok-comment: #94a3b8;
    --light-tok-string: #d97706;
    --light-tok-keyword: #8b5cf6;
    --light-tok-meta: #fb7185;
    --light-tok-number: #059669;
    --light-tok-variable: #e11d48;
    --light-tok-atom: #0284c7;
    --light-tok-def: #e11d48;
  }

  :global(.cm-editor) {
    font-size: 15px;
    font-family: "Source Code Pro", monospace;
    height: 100%;
    width: 100%;
    border-radius: 6px;
    overflow: "hidden";
  }

  :global(.cm-gutter) {
    height: 100%;
  }

  :global(.cm-scroller) {
    height: 100%;
    overflow: auto;
  }

  :global(.cm-cursorLayer .cm-cursor) {
    border-left: 6px solid hsl(var(--p));
  }

  :global(.cm-editor.cm-focused) {
    outline: none;
  }

  /* Dark Theme Custom Style */
  :global([data-theme="dark"] .codemirror) {
    color: var(--dark-tok-text);
  }
  :global([data-theme="dark"] .cm-editor .cm-activeLineGutter) {
    background-color: hsl(var(--b2) / 0.3);
    color: whitesmoke;
  }
  :global([data-theme="dark"] .cm-editor .cm-activeLine) {
    background-color: hsl(var(--b2) / 0.3);
  }
  :global(
      [data-theme="dark"]
        .cm-editor.cm-focused
        .cm-scroller
        .cm-selectionLayer
        .cm-selectionBackground,
      .cm-editor .cm-scroller .cm-selectionLayer .cm-selectionBackground
    ) {
    background-color: hsl(var(--p));
  }

  /* Editor text highlight theme */
  :global([data-theme="dark"] .tok-comment) {
    color: var(--dark-tok-comment);
    font-style: italic;
  }
  :global([data-theme="dark"] .tok-string) {
    color: var(--dark-tok-string);
  }
  :global([data-theme="dark"] .tok-keyword) {
    color: var(--dark-tok-keyword);
  }
  :global([data-theme="dark"] .tok-meta) {
    color: var(--dark-tok-meta);
  }
  :global([data-theme="dark"] .tok-number) {
    color: var(--dark-tok-number);
  }
  :global([data-theme="dark"] .tok-variable) {
    color: var(--dark-tok-variable);
  }
  :global([data-theme="dark"] .tok-atom) {
    color: var(--dark-tok-atom);
  }
  :global([data-theme="dark"] .tok-def) {
    color: var(--dark-tok-def);
  }

  /* Light Theme Custom Style */
  :global([data-theme="light"] .codemirror) {
    color: var(--light-tok-text);
  }
  :global([data-theme="light"] .cm-editor .cm-activeLineGutter) {
    background-color: hsl(var(--b2) / 0.5);
    color: whitesmoke;
  }
  :global([data-theme="light"] .cm-editor .cm-activeLine) {
    background-color: hsl(var(--b2) / 0.5);
  }
  :global(
      [data-theme="light"]
        .cm-editor.cm-focused
        .cm-scroller
        .cm-selectionLayer
        .cm-selectionBackground,
      .cm-editor .cm-scroller .cm-selectionLayer .cm-selectionBackground
    ) {
    background-color: hsl(var(--p));
  }

  /* Editor text highlight theme */
  :global([data-theme="light"] .tok-comment) {
    color: var(--light-tok-comment);
    font-style: italic;
  }
  :global([data-theme="light"] .tok-string) {
    color: var(--light-tok-string);
  }
  :global([data-theme="light"] .tok-keyword) {
    color: var(--light-tok-keyword);
  }
  :global([data-theme="light"] .tok-meta) {
    color: var(--light-tok-meta);
  }
  :global([data-theme="light"] .tok-number) {
    color: var(--light-tok-number);
  }
  :global([data-theme="light"] .tok-variable) {
    color: var(--light-tok-variable);
  }
  :global([data-theme="light"] .tok-atom) {
    color: var(--light-tok-atom);
  }
  :global([data-theme="light"] .tok-def) {
    color: var(--light-tok-def);
  }
</style>
