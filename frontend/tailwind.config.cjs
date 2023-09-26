import { colorScheme } from "./src/lib/util";

/** @type {import('tailwindcss').Config}*/
const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  plugins: [require("daisyui")],

  daisyui: {
    styled: true,
    themes: [
      {
        light: colorScheme.light,
        dark: colorScheme.dark,
      },
    ],
    base: true,
    utils: true,
    logs: false,
    rtl: false,
    prefix: "",
  },

  theme: {
    extend: {
      transitionTimingFunction: {
        "in-out-cubic": "cubic-bezier(0.65, 0, 0.35, 1)",
      },
    },
  },
};

module.exports = config;
