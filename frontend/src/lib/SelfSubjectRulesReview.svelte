<script lang="ts">
  import type { V1ResourceRule } from "@kubernetes/client-node";
  import EmbeddedTable from "../lib/tables/EmbeddedTable.svelte";

  export let rules: Array<V1ResourceRule>;
  export let tagName = "";
</script>

<EmbeddedTable
  tableType={"custom-body"}
  tableItems={[{ name: "ApiGroups" }, { name: "Resources" }, { name: "Verbs" }]}
  {tagName}
>
  <tbody class="text-sm">
    {#if rules}
      {#each rules as rule}
        <tr>
          <td
            class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
          >
            {JSON.stringify(rule.apiGroups, null, 2) ?? "-"}</td
          >
          <td
            class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
          >
            {JSON.stringify(rule.resources, null, 2) ?? "-"}</td
          >
          <td
            class="max-h-fit max-w-fit rounded-lg border border-base-300 bg-base-200/10 px-1 font-mono text-sm"
          >
            {JSON.stringify(rule.verbs, null, 2) ?? "-"}</td
          >
        </tr>
      {/each}
    {/if}
  </tbody>
</EmbeddedTable>
