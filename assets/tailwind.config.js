/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors')

module.exports = {
  content: [
    "../*.{go,html}",
    "../assets/*.{js,ts}"
  ],
  plugins: [],
  mode: 'jit',
  darkMode: 'media',
  theme: {
    extend: {
      colors: {
        gray: colors.neutral
      }
    }
  }
}
