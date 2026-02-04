# Styling

## Overview

RoundTiming uses [TailwindCSS](https://tailwindcss.com/) for styling with utility-first CSS classes. The project also includes the `@tailwindcss/forms` plugin for better form styling.

## Configuration

### tailwind.config.js

```javascript
module.exports = {
  content: [
    './**/*.templ',  // Scan all templ files for classes
  ],
  darkMode: 'class',  // Enable dark mode via class
  plugins: [
    require('@tailwindcss/forms'),  // Form styling plugin
  ],
  corePlugins: {
    preflight: true,  // Enable base styles
  }
}
```

### Input File

`input.css` contains:

- Tailwind directives
- Base layer customizations
- Custom utility classes

## Base Styles

Custom base styles in `input.css`:

### Form Elements

```css
@layer base {
  [type="text"],
  [type="password"],
  [type="email"],
  /* ... other input types ... */
  select,
  textarea {
    @apply shadow-sm rounded-md block w-full sm:text-sm border-gray-300
        focus:ring-orange-400 focus:border-orange-400 focus:outline-none focus:ring-inset
        dark:focus:border-orange-400 dark:bg-gray-900 dark:border-gray-600 dark:text-gray-200;
  }

  [type="checkbox"],
  [type="radio"] {
    @apply rounded focus:ring-orange-400 dark:bg-gray-900;
  }
}
```

### Typography

```css
@layer base {
  h1 {
    @apply text-2xl;
  }
  h2 {
    @apply text-xl;
  }
  h3 {
    @apply text-lg;
  }
}
```

### Links

```css
@layer base {
  .content a,
  .breadcrumbs a {
    @apply text-green-500 hover:text-green-700;
  }
  footer a {
    @apply text-gray-500 hover:text-gray-700 font-bold;
  }
}
```

## Custom Classes

### Tooltip

```css
.tooltip {
  @apply invisible absolute;
}

.has-tooltip:hover .tooltip {
  @apply visible z-50;
}
```

Usage:

```go
<div class="has-tooltip">
    Hover me
    <span class="tooltip">Tooltip text</span>
</div>
```

### HTMX Animations

```css
tr.htmx-swapping td {
  opacity: 0;
  transition: opacity 1s ease-out;
}
```

## Color Palette

The project uses Tailwind's default colors with emphasis on:

| Usage | Color |
| ----- | ----- |
| Primary | `sky-600` |
| Success | `green-400/500` |
| Error | `red-500/600` |
| Warning | `yellow-500` |
| Info | `sky-500` |
| Focus rings | `orange-400` |
| Text | `gray-500/700` |
| Dark mode bg | `gray-900` |

## Dark Mode

Dark mode is enabled via the `class` strategy. Add `dark` class to `<html>`:

```html
<html class="dark">
```

Dark variants in classes:

```go
<div class="bg-white dark:bg-gray-900 text-gray-700 dark:text-gray-200">
```

## Building CSS

### Development (Watch Mode)

```bash
make live/tailwind
# or
npx tailwindcss -i input.css -o public/tailwind.css --watch
```

### Production (Minified)

```bash
make build/tailwind
# or
npx tailwindcss -i input.css -o public/tailwind.css --minify
```

## Using Classes in Templ

### Static Classes

```go
<div class="container mx-auto p-4">
```

### Dynamic Classes

```go
<div class={
    "base-class",
    templ.KV("active", isActive),
    templ.KV("disabled", isDisabled),
}>
```

### Conditional Classes Example

```go
templ PopinMessage(popinMessages PopinMessages) {
    <div class={
        "alert p-4 rounded",
        templ.KV("bg-green-50 border-green-500", popinMessages.Type == "success"),
        templ.KV("bg-red-50 border-red-500", popinMessages.Type == "error"),
        templ.KV("bg-yellow-50 border-yellow-500", popinMessages.Type == "warning"),
    }>
        { children... }
    </div>
}
```

## Best Practices

1. **Use Tailwind utilities** - Avoid writing custom CSS when possible
2. **Extract components** - Use `@apply` for repeated patterns
3. **Mobile-first** - Start with mobile styles, add `sm:`, `md:`, `lg:` breakpoints
4. **Dark mode support** - Always include `dark:` variants for key elements
5. **Cache busting** - CSS includes timestamp for cache invalidation:

```go
href={ fmt.Sprintf("/public/tailwind.css?build=%s", strconv.FormatInt(time.Now().Unix(), 10)) }
```
