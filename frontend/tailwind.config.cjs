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
};

module.exports = config;
