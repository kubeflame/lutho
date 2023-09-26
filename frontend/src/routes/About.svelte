<script lang="ts">
  import CodeMirror from "../lib/CodeMirror.svelte";
  import { icons, tabs } from "../lib/util";

  export let params: any = {};
  params;

  let store: any;

  import { fade } from "svelte/transition";
  import { TabIndex } from "../lib/types";
  import ResourceToolbar from "../lib/ResourceToolbar.svelte";
  import Events from "../lib/Events.svelte";

  // let termParams = { name: "nginx", namespace: "default" };

  let toolbarContent = [
    { index: 0, name: "Root", url: "/about" },
    { index: 1, name: "First", url: "/about/1" },
    { index: 2, name: "Second", url: "/about/2" },
    { index: 3, name: "Third" },
  ];

  import EmbeddedOptions from "../lib/EmbeddedOptions.svelte";

  // import { podList, yaml, kubeEvents, podJson } from "../../../../testData";
  import HeaderElement from "../lib/HeaderElement.svelte";

  const getPodAge = (t: string) => {
    return Date.now() - Date.parse(t);
  };

  import Tabs from "../lib/Tabs.svelte";
  import NamespaceSelect from "../lib/NamespaceSelect.svelte";
  import { availableContainers } from "../lib/stores";

  // $availableContainers = podJson.spec.containers.map((c) => {
  //   return c.name;
  // });

  const tabItems = [
    tabs.details,
    tabs.yaml,
    tabs.logs,
    tabs.shell,
    tabs.events,
    tabs.queue,
    tabs.stack,
  ];

  let activeTab: number = TabIndex.DETAILS;
</script>

<HeaderElement>
  <NamespaceSelect slot="namespace" />
  <Tabs slot="tabs" bind:activeTab {tabItems} />
  <ResourceToolbar slot="toolbar" bind:toolbarContent bind:activeTab />
</HeaderElement>

<div class="router-page" in:fade|global={{ duration: 250 }}>
  {#if activeTab === TabIndex.DETAILS}
    <!-- <div
      class="h-full w-full grid grid-cols-2 p-2 gap-2 overflow-y-scroll overflow-x-scroll"
    > -->
    <!-- <Details {params} /> -->
    <!-- </div> -->
  {:else if activeTab === TabIndex.YAML}
    <!-- <CodeMirror doc={yaml} bind:docStore={store} /> -->
  {:else if activeTab === TabIndex.SHELL}
    <!-- <Shell {termParams} /> -->
  {:else if activeTab === TabIndex.EVENTS}
    <!-- <Events {kubeEvents} /> -->
  {:else if activeTab === TabIndex.QUEUE}
    <!-- <table class="table table-pin-rows">
      <thead>
        <tr class="bg-base-200 shadow-sm">
          <th>Name</th>
          <th>Namespace</th>
          <th>Status</th>
          <th>Age</th>
          <th />
        </tr>
      </thead>
      <tbody>
        {#each podList.items as pod}
          <tr class="">
            <td class="flex gap-x-2">
              <div
                class="badge badge-primary badge-xs mt-0.5 place-self-center"
              />
              <a class="hover:text-primary" href="#/about">
                {pod.metadata.name}
              </a>
            </td>
            <td>{pod.metadata.namespace}</td>
            <td>{pod.status.phase}</td>
            <td>{getPodAge(pod.status.startTime)}</td>
            <td class="flex place-items-center items-center">
              <EmbeddedOptions
                embeddedOptionsData={[
                  {
                    fn: () => {
                      console.log("delete");
                    },
                    classes: "hover:btn-error",
                    icon: icons.trash,
                  },
                  {
                    fn: () => {
                      console.log("exec");
                    },
                    classes: "hover:btn-warning",
                    icon: icons.cmdLine,
                  },
                ]}
              />
            </td>
          </tr>
        {/each}
      </tbody>
    </table> -->
  {:else if activeTab === TabIndex.STACK}
    <!--  -->
  {/if}
</div>
