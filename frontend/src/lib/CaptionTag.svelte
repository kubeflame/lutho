<script lang="ts">
  export let tagName: string;
  export let bgClassNames: string = "bg-base-200/10";
  export let textClassNames: string = "";
</script>

{#if tagName}
  <caption
    class="badge badge-outline mb-3 flex rounded-md border-base-300 font-mono text-sm
    font-thin tracking-wide text-base-content/60 {bgClassNames} {textClassNames}"
  >
    {tagName}
  </caption>
{/if}
