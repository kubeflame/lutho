<script lang="ts">
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import SidebarSlider from "./lib/SidebarSlider.svelte";
  import { sidebarItems } from "./lib/util";
  import { kubeAuthState } from "./lib/stores";
  import AuthPage from "./lib/AuthPage.svelte";
  import { onMount } from "svelte";

  onMount(() => {
    const htmlEl = document.querySelector("html");
    htmlEl?.setAttribute(
      "data-theme",
      (localStorage.getItem("theme") as string) ?? "light",
    );
  });
</script>

{#if $kubeAuthState === "true"}
  <SidebarSlider {sidebarItems} />
  <Router {routes} />
{:else}
  <AuthPage />
{/if}
