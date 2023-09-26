<script lang="ts">
  import "@xterm/xterm/css/xterm.css";
  import BottomAlert from "../BottomAlert.svelte";
  import LoadingNewton from "../LoadingNewton.svelte";
  import createXterm from "./xterm";
  import { transitionEffects } from "../util";
  import { type Writable } from "svelte/store";

  export let name: any;
  export let namespace: any;
  export let showShellReconnect: Writable<boolean>;
  export let activeContainer: Writable<string>;

  let el: HTMLElement;

  const { element, isLoading, alert } = createXterm(
    name,
    namespace,
    activeContainer,
    showShellReconnect,
  );

  $: $element = el;
</script>

{#if $isLoading}
  <LoadingNewton />
{/if}

<div bind:this={el} class="terminal h-full w-full {transitionEffects}">
  <BottomAlert bind:alert={$alert} />
</div>
