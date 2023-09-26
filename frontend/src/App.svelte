<script lang="ts">
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import SidebarSlider from "./lib/SidebarSlider.svelte";
  import { onMount } from "svelte";

  onMount(() => {
    const htmlEl = document.querySelector("html");
    htmlEl?.setAttribute(
      "data-theme",
      (localStorage.getItem("theme") as string) ?? "light",
    );
  });
</script>

<SidebarSlider />
<Router {routes} />
