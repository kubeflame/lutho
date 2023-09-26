<script lang="ts">
  import type { EventsV1Event } from "@kubernetes/client-node/dist/gen/api";
  import { EventType } from "./types";
  import { icons } from "./util";
  import EmbeddedTable from "./EmbeddedTable.svelte";

  export let kubeEvents: EventsV1Event[];
  export let eventsReversed: number = 0;

  if (eventsReversed === 0) kubeEvents.reverse(), (eventsReversed += 1);
</script>

<ol class="m-3 ml-5">
  {#each kubeEvents as event}
    <li
      class="border-l-2
        {event.type === EventType.NORMAL ? 'border-primary' : 'border-warning'}"
    >
      <div class="flex-start md:flex">
        <div
          class="-ml-3.5 flex h-6 w-6 items-center justify-center rounded-full border-2 bg-base-100
            {event.type === EventType.NORMAL
            ? 'border-primary'
            : 'border-warning'}"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width={1.5}
            stroke="currentColor"
            class="h-4 w-4 text-base-content"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d={icons.clock}
            />
          </svg>
        </div>
        <div class="mb-8 ml-4 mt-0.5 w-full">
          <p class="mb-2 font-mono text-sm text-base-content/70">
            {new Date(event.metadata?.creationTimestamp ?? "")}
          </p>
          <div
            class="block rounded-md bg-base-100 pb-4 pl-6 pr-6 pt-2 shadow-md"
          >
            <EmbeddedTable
              embeddedTableItems={[
                { name: "EventType", value: event.type },
                { name: "Reason", value: event.reason },
                {
                  name: "Resource",
                  value: event.regarding?.kind,
                },
                {
                  name: "Name",
                  value: event.regarding?.name,
                },
                {
                  name: "Namespace",
                  value: event.regarding?.namespace,
                  hidden: !event.regarding?.namespace,
                },
                {
                  name: "Message",
                  value: event.note ?? "-",
                },
              ]}
            />
          </div>
        </div>
      </div>
    </li>
  {/each}
</ol>
