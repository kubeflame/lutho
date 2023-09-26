<script lang="ts">
  // import { fetchResource } from "./kubernetes";
  // import { apiURL } from "./util";
  // import { PodV1GVRK, MetricsV1GV } from "./util";
  // import ErrorPage from "./ErrorPage.svelte";

  // export let params: any;

  // const fetchMetricsParams = {
  //     url: apiURL.get,
  //     namespace: params.namespace,
  //     name: params.name,
  //     resource: "pods",
  //     ...MetricsV1GV,
  // };

  // const fetchPodParams = {
  //     url: apiURL.get,
  //     namespace: params.namespace,
  //     name: params.name,
  //     ...PodV1GVRK,
  // };

  // import {
  //     podList,
  //     yaml,
  //     kubeEvents,
  //     podJson,
  //     podMetrics,
  // } from "../../../../testData";
  // import EmbeddedTable from "./EmbeddedTable.svelte";
</script>

<div
  class="grid h-full w-full grid-cols-2 gap-2 overflow-x-scroll overflow-y-scroll p-2"
>
  <!-- <div class="card card-compact flex grow bg-base-100 shadow-md rounded-md">
        <div class="card-title justify-center">Pod Info</div>
        <div class="card-body">
            <EmbeddedTable
                embeddedTableItems={[
                    { name: "Pod", value: podJson.metadata.name },
                    { name: "Namespace", value: podJson.metadata.namespace },
                    {
                        name: "CreationTimestamp",
                        value: podJson.metadata.creationTimestamp,
                    },
                ]}
            />
            {#each podJson.spec.containers as container}
                <EmbeddedTable
                    embeddedTableItems={[
                        { name: "Container", value: container.name },
                        { name: "Image", value: container.image },
                    ]}
                />
            {/each}
        </div>
    </div>
    <div class="card card-compact flex grow bg-base-100 shadow-md rounded-md"> -->
  <!-- <div class="card-title justify-center">Pod Info</div> -->
  <!-- <div class="card-body">
            {#each podMetrics.containers as container}
                <EmbeddedTable
                    embeddedTableItems={[
                        { name: "Container", value: container.name },
                        { name: "CPU usage", value: container.usage.cpu },
                        { name: "Memory usage", value: container.usage.memory },
                    ]}
                />
            {/each} -->
  <!-- <slot name="podInfo" /> -->
  <!-- <h2 class="card-title">Life hack</h2>
        <p>
            How to park your car at your
            garagesajdhgjhsagdhjagshjdgjhagdhjagjshdghja ashgdhjaghjdgjahgs?
            asdjgasjhdgahjsgdhjasghjdaghdjgajhsdghjasgdjhasgdhjgashjdgajhdgajshgdas
        </p>
        <div class="card-actions justify-end">
            <button class="btn btn-primary">Learn now!</button>
        </div> -->
  <!-- </div>
    </div> -->
</div>
<!-- <div>
    {#await fetchResource(fetchMetricsParams) then { error, body }}
        {#if error}
            <ErrorPage errorMessage={error} />
        {:else if body.containers}
            {#each body.containers as container}
                Container Name: {container.name} <br />
                Usage: <br />
                CPU: {container.usage.cpu} <br />
                Memory: {container.usage.memory} <br />
            {/each}
        {/if}
    {/await}

    {#await fetchResource(fetchPodParams) then { error, body }}
        {#if error}
            <ErrorPage errorMessage={error} />
        {:else}
            Name: {body.metadata.name} <br />
            Namespace: {body.metadata.namespace} <br />
            CreationTimestamp: {body.metadata.creationTimestamp} <br />
            NodeName: {body.spec.nodeName} <br />
            Containers: <br />
            {#each body.spec.containers as container}
                Name: {container.name}<br />
                Image: {container.image}<br />
            {/each}
            Conditions:<br />
            {#each body.status.conditions as condition}
                Type: {condition.type}<br />
                Status: {condition.status}<br />
                LastProbeTime: {condition.lastProbeTime}<br />
                LastTransitionTime: {condition.lastTransitionTime}<br />
            {/each}
            Volumes:<br />
            {#if body.spec.volumes}
                {#each body.spec.volumes as volume}
                    Name: {volume.name}<br />
                {/each}
            {/if}
        {/if}
    {/await}
</div> -->
