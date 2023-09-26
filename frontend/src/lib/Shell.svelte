<script lang="ts">
  import "xterm/css/xterm.css";
  import { availableContainers } from "./stores";
  import BottomAlert from "./BottomAlert.svelte";
  import LoadingNewton from "./LoadingNewton.svelte";
  import createXterm from "./xterm";

  export let termParams: any;
  export let showShellReconnect: boolean;
  export let activeContainer: string = "";

  let el: HTMLElement;

  const { container, element, isLoading, alert, shellReconnect } = createXterm({
    name: termParams.name,
    namespace: termParams.namespace,
  });

  $: if (!termParams.shellEmbed) {
    $availableContainers = termParams.containers;
  }
  $: if (activeContainer === "") {
    activeContainer = termParams.containers[0] ?? $availableContainers[0];
  }

  $: $container = activeContainer;
  $: showShellReconnect = $shellReconnect;
  $: $element = el;
</script>

{#if $isLoading}
  <LoadingNewton />
{/if}

<div
  bind:this={el}
  class="terminal h-full w-full transform duration-300 ease-in-out"
>
  <BottomAlert bind:alert={$alert} />
</div>
