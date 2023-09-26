<script lang="ts">
  import { fade } from "svelte/transition";
  import { link } from "svelte-spa-router";
  import { httpStatus } from "./util";
  import { statusCode } from "./stores";

  export let errorMessage: string =
    "Whoops, something went wrong on our servers.";
</script>

<div
  class="error-page flex h-full flex-col items-center justify-center space-x-8 space-y-16 lg:flex-row lg:space-y-0 2xl:space-x-0"
  in:fade|global={{ duration: 250 }}
>
  <div
    class="flex w-full flex-col items-center justify-center text-center lg:w-1/2 lg:px-2 xl:px-0"
  >
    <p
      class="text-7xl font-bold tracking-wider text-error drop-shadow-md md:text-8xl lg:text-9xl"
    >
      {$statusCode ?? 503}
    </p>
    <p
      class="text-4xl font-bold tracking-wider text-error drop-shadow-sm md:text-5xl lg:text-6xl"
    >
      {httpStatus[$statusCode] ?? "Service Unavailable"}
    </p>
    <p class="my-12 text-lg md:text-xl lg:text-2xl">
      {errorMessage}
    </p>

    <div
      class="flex w-full flex-row items-center justify-center gap-4 text-center"
    >
      <button
        class="btn btn-ghost flex bg-base-200 hover:bg-primary hover:drop-shadow-md"
        on:click={() => window.location.reload()}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="h-5 w-5 pr-1"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99"
          />
        </svg>
        Refresh Page
      </button>
      <a
        role="button"
        class="btn btn-ghost flex bg-base-200 transition duration-150 hover:bg-primary hover:drop-shadow-md"
        href="/"
        use:link
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 pr-1"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M9.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L7.414 9H15a1 1 0 110 2H7.414l2.293 2.293a1 1 0 010 1.414z"
          />
        </svg>
        Return Home
      </a>
    </div>
  </div>

  <div
    class="flex w-1/2 flex-col justify-center pb-20 pr-20 drop-shadow-md lg:h-full lg:items-end"
  >
    <svg
      class="w-full text-error"
      xmlns="http://www.w3.org/2000/svg"
      data-name="Layer 1"
      viewBox="0 0 1119.60911 699"
      xmlns:xlink="http://www.w3.org/1999/xlink"
    >
      <path
        d="M31.39089,151.64237c0,77.49789,48.6181,140.20819,108.70073,140.20819"
        transform="translate(-31.39089 -100.5)"
        fill="hsl(var(--b3))"
      />
      <path
        d="M140.09162,291.85056c0-78.36865,54.255-141.78356,121.30372-141.78356"
        transform="translate(-31.39089 -100.5)"
        fill="currentColor"
      />
      <path
        d="M70.77521,158.66768c0,73.61476,31.00285,133.18288,69.31641,133.18288"
        transform="translate(-31.39089 -100.5)"
        fill="currentColor"
      />
      <path
        d="M140.09162,291.85056c0-100.13772,62.7103-181.16788,140.20819-181.16788"
        transform="translate(-31.39089 -100.5)"
        fill="hsl(var(--b3))"
      />
      <ellipse
        cx="198.60911"
        cy="424.5"
        rx="187"
        ry="25.43993"
        fill="hsl(var(--b3))"
        class="drop-shadow-md"
      />
      <ellipse
        cx="836.60911"
        cy="660.5"
        rx="283"
        ry="38.5"
        fill="hsl(var(--b3))"
        class="drop-shadow-md"
      />
      <ellipse
        cx="310.60911"
        cy="645.5"
        rx="170"
        ry="23.12721"
        fill="hsl(var(--b3)/0.8)"
        class="drop-shadow-md"
      />
      <path
        d="M494,726.5c90,23,263-30,282-90"
        transform="translate(-31.39089 -100.5)"
        fill="none"
        stroke="hsl(var(--b3))"
        stroke-miterlimit="10"
        stroke-width="2"
      />
      <path
        d="M341,359.5s130-36,138,80-107,149-17,172"
        transform="translate(-31.39089 -100.5)"
        fill="none"
        stroke="hsl(var(--b3))"
        stroke-miterlimit="10"
        stroke-width="2"
      />
      <path
        d="M215.40233,637.78332s39.0723-10.82,41.47675,24.04449-32.15951,44.78287-5.10946,51.69566"
        transform="translate(-31.39089 -100.5)"
        fill="none"
        stroke="hsl(var(--b3))"
        stroke-miterlimit="10"
        stroke-width="2"
      />
      <path
        d="M810.09554,663.73988,802.218,714.03505s-38.78182,20.60284-11.51335,21.20881,155.73324,0,155.73324,0,24.84461,0-14.54318-21.81478l-7.87756-52.719Z"
        transform="translate(-31.39089 -100.5)"
        fill="hsl(var(--b3))"
      />
      <rect
        x="578.43291"
        y="212.68859"
        width="513.25314"
        height="357.51989"
        rx="18.04568"
        fill="hsl(var(--b3))"
      />
      <rect
        x="595.70294"
        y="231.77652"
        width="478.71308"
        height="267.83694"
        fill="hsl(var(--b1))"
      />
      <circle cx="835.05948" cy="223.29299" r="3.02983" fill="hsl(var(--b1))" />
      <polygon
        points="968.978 667.466 968.978 673.526 642.968 673.526 642.968 668.678 643.417 667.466 651.452 645.651 962.312 645.651 968.978 667.466"
        fill="hsl(var(--b1)/0.8)"
        class="shadow-lg shadow-inherit"
      />
      <polygon
        points="968.978 667.466 968.978 673.526 642.968 673.526 642.968 668.678 643.417 667.466 968.978 667.466"
        fill="hsl(var(--b2))"
        class="shadow-lg shadow-inherit"
      />
      <path
        d="M1125.828,762.03359c-.59383,2.539-2.83591,5.21743-7.90178,7.75032-18.179,9.08949-55.1429-2.42386-55.1429-2.42386s-28.4804-4.84773-28.4804-17.573a22.72457,22.72457,0,0,1,2.49658-1.48459c7.64294-4.04351,32.98449-14.02122,77.9177.42248a18.73921,18.73921,0,0,1,8.54106,5.59715C1125.07908,756.45353,1126.50669,759.15715,1125.828,762.03359Z"
        transform="translate(-31.39089 -100.5)"
        fill="hsl(var(--b1)/0.6)"
        class="shadow-lg shadow-inherit"
      />
      <path
        d="M1125.828,762.03359c-22.251,8.526-42.0843,9.1622-62.43871-4.975-10.26507-7.12617-19.59089-8.88955-26.58979-8.75618,7.64294-4.04351,32.98449-14.02122,77.9177.42248a18.73921,18.73921,0,0,1,8.54106,5.59715C1125.07908,756.45353,1126.50669,759.15715,1125.828,762.03359Z"
        transform="translate(-31.39089 -100.5)"
        fill="hsl(var(--b1)/0.8)"
        class="shadow-lg shadow-inherit"
      />
      <circle
        cx="835.05948"
        cy="545.66686"
        r="11.51335"
        fill="hsl(var(--b1))"
      />
      <rect
        x="108.60911"
        y="159"
        width="208"
        height="242"
        fill="hsl(var(--b3)/0.9)"
        class="drop-shadow-md"
      />
      <rect
        x="87.60911"
        y="135"
        width="250"
        height="86"
        fill="hsl(var(--b2))"
        class="drop-shadow-md"
      />
      <rect
        x="87.60911"
        y="237"
        width="250"
        height="86"
        fill="hsl(var(--er))"
        class="drop-shadow-md"
      />
      <rect
        x="87.60911"
        y="339"
        width="250"
        height="86"
        fill="hsl(var(--b2))"
        class="drop-shadow-md"
      />
      <rect
        x="271.60911"
        y="150"
        width="16"
        height="16"
        fill="hsl(var(--er)/0.3)"
      />
      <rect
        x="294.60911"
        y="150"
        width="16"
        height="16"
        fill="hsl(var(--er)/0.6)"
      />
      <rect
        x="317.60911"
        y="150"
        width="16"
        height="16"
        fill="hsl(var(--er)/0.9)"
      />
      <rect
        x="271.60911"
        y="251"
        width="16"
        height="16"
        fill="hsl(var(--wa)/0.3)"
        class="drop-shadow-sm"
      />
      <rect
        x="294.60911"
        y="251"
        width="16"
        height="16"
        fill="hsl(var(--wa)/0.6)"
        class="drop-shadow-sm"
      />
      <rect
        x="317.60911"
        y="251"
        width="16"
        height="16"
        fill="hsl(var(--wa)/0.9)"
        class="drop-shadow-sm"
      />
      <rect
        x="271.60911"
        y="352"
        width="16"
        height="16"
        fill="currentColor"
        opacity="0.4"
      />
      <rect
        x="294.60911"
        y="352"
        width="16"
        height="16"
        fill="currentColor"
        opacity="0.8"
      />
      <rect x="317.60911" y="352" width="16" height="16" fill="currentColor" />
      <circle cx="316.60911" cy="538" r="79" fill="hsl(var(--er)/0.9)" />
      <rect
        x="280.60911"
        y="600"
        width="24"
        height="43"
        fill="hsl(var(--er)/0.9)"
      />
      <rect
        x="328.60911"
        y="600"
        width="24"
        height="43"
        fill="hsl(var(--er)/0.9)"
      />
      <ellipse
        cx="300.60911"
        cy="643.5"
        rx="20"
        ry="7.5"
        fill="hsl(var(--er)/0.9)"
      />
      <ellipse
        cx="348.60911"
        cy="642.5"
        rx="20"
        ry="7.5"
        fill="hsl(var(--er)/0.9)"
      />
      <circle cx="318.60911" cy="518" r="27" fill="#fff" />
      <circle cx="318.60911" cy="518" r="9" fill="hsl(var(--er))" />
      <path
        d="M271.36733,565.03228c-6.37889-28.56758,14.01185-57.43392,45.544-64.47477s62.2651,10.41,68.644,38.9776-14.51861,39.10379-46.05075,46.14464S277.74622,593.59986,271.36733,565.03228Z"
        transform="translate(-31.39089 -100.5)"
        fill="hsl(var(--b3))"
      />
      <ellipse
        cx="417.21511"
        cy="611.34365"
        rx="39.5"
        ry="12.40027"
        transform="translate(-238.28665 112.98044) rotate(-23.17116)"
        fill="hsl(var(--er)/0.9)"
      />
      <ellipse
        cx="269.21511"
        cy="664.34365"
        rx="39.5"
        ry="12.40027"
        transform="translate(-271.07969 59.02084) rotate(-23.17116)"
        fill="hsl(var(--er)/0.9)"
      />
      <path
        d="M394,661.5c0,7.732-19.90861,23-42,23s-43-14.268-43-22,20.90861-6,43-6S394,653.768,394,661.5Z"
        transform="translate(-31.39089 -100.5)"
        fill="#fff"
      />
    </svg>
  </div>
</div>
