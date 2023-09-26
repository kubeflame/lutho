<script lang="ts">
  import { icons } from "./util";
  import { showSettings, sidebarState } from "./stores";
  import { slide } from "svelte/transition";
  import { cubicInOut } from "svelte/easing";
  import Settings from "./Settings.svelte";

  const htmlEl = document.querySelector("html");
  const logoName = "LUTHO";
  const fontLogo = "font-mono font-thin tracking-wider";

  let darkMode: boolean = localStorage.getItem("theme") === "dark";

  function setDataTheme() {
    if (darkMode) {
      htmlEl?.setAttribute("data-theme", "dark");
      localStorage.setItem("theme", "dark");
    } else {
      htmlEl?.setAttribute("data-theme", "light");
      localStorage.setItem("theme", "light");
    }
  }

  function setState() {
    $sidebarState = $sidebarState === "max" ? "min" : "max";
    localStorage.setItem("navState", $sidebarState);
  }
</script>

<Settings />

<div class="nav-button z-10 flex">
  <button
    class="nav-state-button z-50 flex self-center rounded-full
      border-2 border-base-100 bg-base-200 p-0.5 duration-300 ease-in-out
      hover:rotate-90 hover:bg-primary hover:drop-shadow-md"
    on:click={setState}
  >
    <svg
      width="32px"
      height="32px"
      viewBox="0 0 64 64"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill="#274c77"
        d="M46.2,32c0-.9,.3-1.72,.8-2.4-.36-2.14-1.2-4.11-2.4-5.81-.83-.13-1.63-.49-2.26-1.13s-1.01-1.44-1.13-2.26c-1.7-1.2-3.67-2.04-5.81-2.4-.67,.5-1.5,.8-2.4,.8s-1.72-.3-2.4-.8c-2.14,.36-4.11,1.2-5.81,2.4-.13,.83-.49,1.63-1.13,2.26s-1.44,1.01-2.26,1.13c-1.2,1.7-2.04,3.67-2.4,5.81,.5,.67,.8,1.5,.8,2.4s-.3,1.72-.8,2.4c.36,2.14,1.2,4.11,2.4,5.81,.83,.13,1.63,.49,2.26,1.13,.64,.64,1.01,1.44,1.13,2.26,1.7,1.2,3.67,2.04,5.81,2.4,.67-.5,1.5-.8,2.4-.8s1.72,.3,2.4,.8c2.14-.36,4.11-1.2,5.81-2.4,.13-.83,.49-1.63,1.13-2.26,.64-.64,1.44-1.01,2.26-1.13,1.2-1.7,2.04-3.67,2.4-5.81-.5-.67-.8-1.5-.8-2.4Z"
      />
      <circle cx="33" cy="32" r="9" fill="#f5dd90" />
      <path
        fill="#6096ba"
        d="M7,32c0-14.36,11.64-26,26-26v4c-12.15,0-22,9.85-22,22"
      />
      <polygon fill="#f76c5e" points="9 38 5 32 13 32 9 38" />
      <path
        fill="#6096ba"
        d="M46,54.51c-12.44,7.18-28.34,2.92-35.52-9.52l3.46-2c6.08,10.52,19.53,14.13,30.05,8.05"
      />
      <polygon
        fill="#f76c5e"
        points="50.19 49.78 47 56.25 43 49.32 50.19 49.78"
      />
      <path
        fill="#6096ba"
        d="M46,9.48c12.44,7.18,16.7,23.08,9.52,35.52l-3.46-2c6.08-10.52,2.47-23.98-8.05-30.05"
      />
      <polygon fill="#f76c5e" points="39.8 8.21 47 7.75 43 14.68 39.8 8.21" />
    </svg>
  </button>

  {#if $sidebarState === "max"}
    <div
      transition:slide={{ axis: "x", duration: 300, easing: cubicInOut }}
      class="nav-toolbar -ml-7 flex h-8 w-56 items-center self-center rounded-full
        border-2 border-base-100 bg-base-200 pl-7 shadow-sm"
    >
      <div class="nav-options flex h-full items-center">
        <label
          class="swap-rotate swap ml-1 rounded-full p-1 hover:bg-primary hover:drop-shadow-md"
        >
          <input
            id="data-toggle-theme"
            type="checkbox"
            bind:checked={darkMode}
            on:change={setDataTheme}
          />
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width={1.5}
            stroke="currentColor"
            class="swap-on h-[1.17rem] w-[1.17rem]"
          >
            <path d={icons.moon} />
          </svg>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width={1.5}
            stroke="currentColor"
            class="swap-off h-[1.17rem] w-[1.17rem]"
          >
            <path d={icons.sun} />
          </svg>
        </label>

        <button
          class="z-50 rounded-full p-1 hover:bg-primary hover:drop-shadow-md"
          on:click={() => {
            $showSettings = true;
          }}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width={1.25}
            stroke="currentColor"
            class="h-[1.17rem] w-[1.17rem]"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d={icons.cogWheel}
            />
          </svg>
        </button>
      </div>

      <div
        class="nav-logo group flex h-full items-center rounded-full bg-gradient-to-r
          from-base-200 via-primary to-primary"
      >
        <div class="group-item flex duration-300 ease-in-out {fontLogo}">
          {logoName}
        </div>
      </div>
    </div>
  {/if}
</div>
