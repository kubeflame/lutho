<script lang="ts">
  import { icons } from "./icons";

  export let type: keyof typeof icons = "hashtag";
  export let classNames: string = "h-4 w-4";
  export let viewBox: string = "0 0 24 24";
  export let strokeWidth: number = 1.5;
  export let stroke: string = "currentColor";
  export let fill: string = "none";
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  {fill}
  {viewBox}
  stroke-width={strokeWidth}
  {stroke}
  class={classNames}
>
  <path stroke-linecap="round" stroke-linejoin="round" d={icons[type]} />
</svg>
