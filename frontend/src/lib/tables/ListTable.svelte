<script lang="ts">
  import { link } from "svelte-spa-router";
  import SvgIcon from "../SvgIcon.svelte";
  import type { EmbeddedOptionsData } from "../types";
  import EmbeddedOptions from "./EmbeddedOptions.svelte";

  export let tableHead: string[];
  export let items: any[];
  export let hrefRoot: string;
  export let isNamespaced: boolean;
  export let embeddedOptionsDataFn: (item: any) => EmbeddedOptionsData[];
  export let hasEmbeddedOptions: boolean = true;
</script>

<table class="table table-pin-rows">
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
        {#if item.metadata.deletionGracePeriodSeconds}
          <tr class="bg-error/20 opacity-70">
            <td class="flex items-center gap-x-2">
              <slot name="podPhaseBadge" {item} />
              {item.metadata?.name}
              <slot name="nodeSchedule" {item} />
            </td>
            {#if isNamespaced}
              <td>{item.metadata?.namespace}</td>
            {/if}
            <slot name="podPhase" {item} />
            <td>{item.metadata?.creationTimestamp}</td>
            <slot name="nodeStatus" {item} />
            <td
              class="flex place-items-center items-center justify-end space-x-5"
            >
              <span
                class="badge badge-outline badge-md border-error text-base-content"
              >
                deletion in progress...
              </span>
              <SvgIcon type="warning" classNames="size-5 pr-1" />
            </td>
          </tr>
        {:else}
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
            <td class="flex place-items-center items-center justify-end">
              {#if hasEmbeddedOptions}
                <EmbeddedOptions
                  embeddedOptionsData={embeddedOptionsDataFn(item)}
                />
              {/if}
            </td>
          </tr>
        {/if}
      {/each}
    {/if}
  </tbody>
</table>
