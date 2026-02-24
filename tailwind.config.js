/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: ['./**/*.templ'],
    extract: {
      templ: (content) => {
        // Default word extraction
        const defaultMatches = content.match(/[^<>"'`\s]*[^<>"'`\s:]/g) || []
        // Extract class values from templ.Attributes{"class": "..."} patterns
        const attrMatches = []
        const classPattern = /"class"\s*:\s*"([^"]+)"/g
        let match
        while ((match = classPattern.exec(content)) !== null) {
          attrMatches.push(...match[1].split(/\s+/))
        }
        return [...defaultMatches, ...attrMatches]
      },
    },
  },
  darkMode: 'class',
  // Safelist for dynamic classes generated at runtime (e.g., team colors)
  safelist: [
    // Team background colors
    { pattern: /bg-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(50|100|200|300|400|500)/ },
    // Team background colors (dark mode)
    { pattern: /bg-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(800|900)/, variants: ['dark'] },
    // Team text colors
    { pattern: /text-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(500|600|700|800)/ },
    // Team border colors
    { pattern: /border-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(200|300|400|500)/ },
    // Team border colors (dark mode)
    { pattern: /border-(red|blue|green|yellow|purple|pink|orange|cyan|indigo|teal)-(700|800)/, variants: ['dark'] },
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
