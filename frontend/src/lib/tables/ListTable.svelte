<script lang="ts">
  import { link } from "svelte-spa-router";

  export let tableHead: string[];
  export let items: any[];
  export let hrefRoot: string;
  export let isNamespaced: boolean;
</script>

<table class="table-pin-rows table">
  <thead>
    <tr class="bg-base-200 shadow-sm">
      {#each tableHead as th}
        <th>{th}</th>
      {/each}
    </tr>
  </thead>
  <tbody>
    {#if items}
      {#each items as item}
        <tr class="">
          <td class="flex items-center gap-x-2">
            <slot name="podPhaseBadge" {item} />
            <a
              class="hover:text-primary"
              href="{hrefRoot}/{isNamespaced
                ? item.metadata?.namespace + '/'
                : ''}{item.metadata?.name}"
              use:link
            >
              {item.metadata?.name}
            </a>
            <slot name="nodeSchedule" {item} />
          </td>
          {#if isNamespaced}
            <td>{item.metadata?.namespace}</td>
          {/if}
          <slot name="podPhase" {item} />
          <td>{item.metadata?.creationTimestamp}</td>
          <slot name="nodeStatus" {item} />
          <slot name="embeddedOptions" {item} />
        </tr>
      {/each}
    {/if}
  </tbody>
</table>
