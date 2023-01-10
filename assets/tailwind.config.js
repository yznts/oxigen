/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');

module.exports = {
  content: [
    "../*.{go,html}",
    "../assets/*.{js,ts}"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('rippleui')
  ],
  mode: 'jit',
  darkMode: 'media'
}
