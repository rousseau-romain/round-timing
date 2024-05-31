/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './**/*.templ',
  ],
  darkMode: 'class',
  plugins: [
    require('@tailwindcss/forms'),
  ],
  corePlugins: {
    preflight: true,
  }
}
