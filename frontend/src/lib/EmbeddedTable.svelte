<script lang="ts">
  import type { EmbeddedTableItem } from "./types";

  export let embeddedTableItems: Array<EmbeddedTableItem> = [];
</script>

<div class="overflow-x-auto">
  <table class="table-auto">
    <tbody class="text-sm">
      {#each embeddedTableItems as item}
        {#if !item.hidden}
          <tr>
            <td
              class="flex items-center whitespace-nowrap pr-4 font-mono font-medium tracking-wide text-base-content"
            >
              {item.name}
            </td>
            <td class="text-left text-base-content/70">
              {item.value}
            </td>
          </tr>
        {/if}
      {/each}
    </tbody>
  </table>
</div>
