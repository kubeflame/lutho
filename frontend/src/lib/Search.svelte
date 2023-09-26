<script lang="ts">
  import { KubeDataOpType, type KubeGVRK } from "./types";
  import socketStore from "./socketStore";

  export let resourceType: KubeGVRK;
  export let namespaceInput: string;
  export let nameInput: string;
  export let labelInput: string;
  export let fieldInput: string;
  export let errorMessage: string;
  export let isFetching: boolean;
  export let onClickSearch: boolean = false;

  let searchDataResponse: any;
  let isNamespaced: boolean;

  const { isLoading, dataSend, sockError, dataList } = socketStore();

  isLoading.subscribe((load) => (isFetching = load));

  $: if (onClickSearch) {
    isNamespaced = false;
    $dataSend = [
      {
        type: KubeDataOpType.list,
        request: {
          namespace: resourceType?.isNamespaced ? namespaceInput : "",
          ...resourceType,
          options: {
            fieldSelector: nameInput
              ? fieldInput.concat(",", `metadata.name=${nameInput}`)
              : fieldInput,
            labelSelector: labelInput,
          },
        },
      },
    ];
    if (resourceType.isNamespaced) isNamespaced = true;
    onClickSearch = false;
  }

  $: errorMessage = $sockError;

  $: searchDataResponse = $dataList;
</script>

{#if searchDataResponse}
  <table class="table table-pin-rows">
    <thead>
      <tr class="bg-base-200 shadow-sm">
        <th>Name</th>
        {#if isNamespaced}
          <th>Namespace</th>
        {/if}
      </tr>
    </thead>
    <tbody>
      {#each searchDataResponse.items as searchItem}
        <tr class="">
          <td class="flex items-center gap-x-2">
            {searchItem.metadata?.name}
          </td>
          {#if isNamespaced}
            <td>{searchItem.metadata?.namespace}</td>
          {/if}
        </tr>
      {/each}
    </tbody>
  </table>
{/if}
