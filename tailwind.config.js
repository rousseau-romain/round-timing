/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './**/*.templ',
  ],
  darkMode: 'class',
  // Safelist for dynamic classes generated at runtime (e.g., team colors)
  safelist: [
    // Team background colors
    { pattern: /bg-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(50|100|200|300|400|500)/ },
    // Team text colors
    { pattern: /text-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(500|600|700|800)/ },
    // Team border colors
    { pattern: /border-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(300|400|500)/ },
  ],
  theme: {
    extend: {
      colors: {
        // Custom semantic colors can be added here
        // primary: colors.sky,
        // success: colors.green,
        // danger: colors.red,
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
  corePlugins: {
    preflight: true,
  }
}
