<script lang="ts">
  import { fade } from "svelte/transition";
  import ErrorPage from "./ErrorPage.svelte";
  import LoadingNewton from "./LoadingNewton.svelte";
  import {
    kubeHost,
    defaultAccessReview,
    sidebarState,
    statusCode,
  } from "./stores";
  import { apiURL, routeString, transitionEffects } from "./util";
  import { onMount } from "svelte";
  import type { AuthResponse } from "./types";
  import { push } from "svelte-spa-router";

  export let errorMessage: string;
  export let loading: boolean;

  onMount(async function () {
    await fetch(apiURL.authState)
      .then((resp) => {
        return resp.json();
      })
      .then((data: AuthResponse) => {
        errorMessage = data.error;
        $kubeHost = data.kubeHost;
        $defaultAccessReview = data.accessReview;
        if (!data.state)
          push(`${routeString.auth}?referrer=${window.location.hash}`);
      })
      .catch((err) => {
        $statusCode = 503;
        errorMessage = err;
        $kubeHost = "";
      });
  });
</script>

<div
  class="content fixed flex h-screen flex-col {transitionEffects}
    {$sidebarState === 'min'
    ? 'ml-12 w-[calc(100vw-48px)]'
    : 'ml-60 w-[calc(100vw-240px)]'}"
>
  <slot name="resource-toolbar" />
  <div
    in:fade
    class="router-page relative m-4 h-[calc(100%-104px)] overflow-y-scroll rounded-md shadow-md"
  >
    {#if errorMessage}
      <ErrorPage bind:errorMessage />
    {:else if loading}
      <LoadingNewton />
    {:else}
      <slot />
    {/if}
  </div>
</div>
