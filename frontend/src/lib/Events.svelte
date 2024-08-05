<script lang="ts">
  import type { CoreV1Event } from "@kubernetes/client-node/dist/gen/api";
  import { EventType } from "./util";
  import EmbeddedTable from "./tables/EmbeddedTable.svelte";
  import SvgIcon from "./SvgIcon.svelte";

  export let kubeEvents: CoreV1Event[];
  export let eventsReversed: number = 0;
  export let svgIconClassNames: string = "h-4 w-4 text-base-content";
  export let svgIconStrokeWidth: number = 1.5;

  $: if (kubeEvents?.length > 0 && eventsReversed === 0)
    kubeEvents.reverse(), (eventsReversed = 1);
</script>

{#if kubeEvents?.length > 0}
  <ol class="m-3 ml-5">
    {#each kubeEvents as event}
      <li
        class="border-l-2
        {event.type === EventType.NORMAL ? 'border-primary' : 'border-warning'}"
      >
        <div class="flex-start md:flex">
          <div
            class="bg-base-100 -ml-3.5 flex h-6 w-6 items-center justify-center rounded-full border-2
            {event.type === EventType.NORMAL
              ? 'border-primary'
              : 'border-warning'}"
          >
            <SvgIcon
              classNames={svgIconClassNames}
              strokeWidth={svgIconStrokeWidth}
              type={"clock"}
            />
          </div>
          <div class="mb-8 ml-4 mt-0.5 w-full">
            <p class="text-base-content/70 mb-2 font-mono text-sm">
              {new Date(event.metadata?.creationTimestamp ?? "")}
            </p>
            <div
              class="bg-base-100 block rounded-md pb-4 pl-6 pr-6 pt-2 shadow-md"
            >
              <EmbeddedTable
                tableType={"events"}
                tableItems={[
                  { name: "EventType", value: event.type },
                  { name: "Reason", value: event.reason },
                  {
                    name: "Resource",
                    value: event.kind,
                  },
                  {
                    name: "Name",
                    value: event.involvedObject?.name,
                  },
                  {
                    name: "Namespace",
                    value: event.involvedObject?.namespace,
                    hidden: !event.involvedObject?.namespace,
                  },
                  {
                    name: "Message",
                    value: event.message ?? "-",
                  },
                ]}
              />
            </div>
          </div>
        </div>
      </li>
    {/each}
  </ol>
{/if}
