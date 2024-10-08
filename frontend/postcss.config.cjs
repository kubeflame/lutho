const tailwindcss = require("tailwindcss");
const autoprefixer = require("autoprefixer");
const cssnano = require("cssnano");

const config = {
  plugins: {
    //Some plugins, like tailwindcss/nesting, need to run before Tailwind,
    tailwindcss: tailwindcss,
    //But others, like autoprefixer, need to run after,
    autoprefixer: autoprefixer,
    cssnano: cssnano,
  },
};

module.exports = config;
