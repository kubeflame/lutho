<script lang="ts">
  import CaptionTag from "../CaptionTag.svelte";
  import type { TableType } from "../types";

  export let tableItems: any;
  export let tableType: TableType;
  export let tagName: string = "";
</script>

<div class="overflow-x-auto">
  {#if tableType === "normal" && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table-auto">
      <thead>
        <tr>
          {#each tableItems.head as head}
            <th>{head}</th>
          {/each}
        </tr>
      </thead>
      <tbody class="text-sm">
        <tr>
          {#each tableItems.body as body}
            <td>{body}</td>
          {/each}
        </tr>
      </tbody>
    </table>
  {:else if tableType === "events" && tableItems.length > 0}
    <table class="table-auto">
      <tbody class="text-sm">
        {#each tableItems as item}
          {#if !item.hidden}
            <tr>
              <td
                class="flex items-center whitespace-nowrap pr-4 font-mono
                  font-medium tracking-wide text-base-content"
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
  {:else if tableType === "one_row" && tableItems && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table table-sm">
      <thead>
        <tr>
          {#each tableItems as el}
            <th class="capitalize {el.className}">{el.name}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each tableItems as el}
          <td class={el.className}>{el.value}</td>
        {/each}
      </tbody>
    </table>
  {:else if tableType === "badges" && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table table-sm">
      {#each tableItems as label}
        <div
          class="badge badge-ghost mx-1 my-3 bg-base-200 text-base-content/70"
        >
          {label.name}: {label.value}
        </div>
      {/each}
    </table>
  {:else if tableType === "custom-body" && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table table-sm">
      <thead>
        <tr class="">
          {#each tableItems as item}
            <th class={item.className}>{item.name}</th>
          {/each}
        </tr>
      </thead>
      <slot />
    </table>
  {:else if tableType === "custom-vertical" && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table table-sm">
      <tbody>
        {#each tableItems as item}
          {#if item.value}
            <tr class="rounded-xl hover:bg-base-200/40">
              <thead>
                <th>{item.name}</th>
              </thead>
              <td>
                <div class={item.className}>
                  {item.value}
                </div>
              </td>
            </tr>
          {/if}
        {/each}
      </tbody>
    </table>
  {/if}
</div>
