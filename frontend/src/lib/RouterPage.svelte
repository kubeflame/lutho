<script lang="ts">
  import { sidebarState } from "./stores";
</script>

<div
  class="content fixed flex h-screen flex-col transition-all duration-300 ease-in-out
  {$sidebarState === 'min'
    ? 'ml-12 w-[calc(100vw-48px)]'
    : 'ml-12 w-[calc(100vw-48px)] md:ml-60 md:w-[calc(100vw-240px)]'}"
>
  <slot name="resource-toolbar" />
  <div
    class="router-page relative m-4 h-[calc(100%-104px)] overflow-y-scroll rounded-md shadow-md"
  >
    <slot />
  </div>
</div>
