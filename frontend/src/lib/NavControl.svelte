<script lang="ts">
  import { showSettings, sidebarState } from "./stores";
  import Settings from "./Settings.svelte";
  import SvgIcon from "./SvgIcon.svelte";
  import { transitionEffects } from "./util";

  const logoName = "LUTHO";
  const fontLogo = "font-mono font-thin tracking-wider";
  const svgIconClassNames: string = "size-5";
  const svgIconStrokeWidth: number = 1.5;

  let darkMode = true;

  function handleSwitchDarkMode() {
    darkMode = !darkMode;

    localStorage.setItem("theme", darkMode ? "dark" : "light");

    darkMode
      ? document.documentElement.setAttribute("data-theme", "dark")
      : document.documentElement.setAttribute("data-theme", "light");
  }

  function setState() {
    $sidebarState = $sidebarState === "max" ? "min" : "max";
    localStorage.setItem("navState", $sidebarState);
  }
</script>

<Settings />

<div class="nav-control no-scrollbar z-10 flex">
  <button
    class="nav-state-button border-base-100 bg-base-200 hover:bg-primary z-50
      flex self-center rounded-full border-2 p-0.5 duration-300 ease-in-out
      hover:rotate-90 hover:drop-shadow-md"
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
      <circle cx="32" cy="32" r="9" fill="#f5dd90" />
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

  <div
    class="nav-toolbar border-base-100 bg-base-200 flex h-8 w-56 items-center
      self-center rounded-full border-2 pl-7 shadow-sm {transitionEffects}
      {$sidebarState === 'max' ? '-ml-7' : '-ml-72'}"
  >
    <div class="nav-options flex h-full w-full items-center">
      <label
        class="swap swap-rotate hover:bg-primary ml-1 rounded-full p-1 hover:drop-shadow-md"
      >
        <input
          class="theme-controller"
          type="checkbox"
          bind:checked={darkMode}
          on:click={handleSwitchDarkMode}
        />
        <SvgIcon
          classNames={`${svgIconClassNames} swap-on`}
          strokeWidth={svgIconStrokeWidth}
          type={"moon"}
        />
        <SvgIcon
          classNames={`${svgIconClassNames} swap-off`}
          strokeWidth={svgIconStrokeWidth}
          type={"sun"}
        />
      </label>

      <button
        class="hover:bg-primary rounded-full p-1 hover:drop-shadow-md"
        on:click={() => {
          $showSettings = true;
        }}
      >
        <SvgIcon
          classNames={svgIconClassNames}
          strokeWidth={svgIconStrokeWidth}
          type={"cogWheel"}
        />
      </button>
    </div>

    <div
      class="nav-logo from-base-200 via-primary to-primary group flex h-full
        w-full items-center rounded-full bg-gradient-to-r pl-7
        {$sidebarState === 'max' ? '' : 'hidden'}"
    >
      <div class="group-item flex {fontLogo}">
        {logoName}
      </div>
    </div>
  </div>
</div>
