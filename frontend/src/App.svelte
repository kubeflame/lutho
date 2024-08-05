<script lang="ts">
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import SidebarSlider from "./lib/SidebarSlider.svelte";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";

  onMount(() => themeChange());
</script>

<SidebarSlider />
<Router {routes} />
