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
                class="text-base-content flex items-center whitespace-nowrap pr-4 font-mono font-medium tracking-wide"
              >
                {item.name}
              </td>
              <td class="text-base-content/70 text-left">
                {item.value}
              </td>
            </tr>
          {/if}
        {/each}
      </tbody>
    </table>
  {:else if tableType === "one_row" && tableItems && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table-sm table">
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
    <table class="table-sm table">
      {#each tableItems as label}
        <div
          class="badge badge-ghost bg-base-200 text-base-content/70 mx-1 my-3"
        >
          {label.name}: {label.value}
        </div>
      {/each}
    </table>
  {:else if tableType === "custom-body" && tableItems.length > 0}
    <CaptionTag {tagName} />
    <table class="table-sm table">
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
    <table class="table-sm table">
      <tbody>
        {#each tableItems as item}
          {#if item.value}
            <tr class="hover:bg-base-200/40 rounded-xl">
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
