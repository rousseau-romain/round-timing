# Translations (i18n)

## Overview

RoundTiming uses [ctxi18n](https://github.com/invopop/ctxi18n) for internationalization. Translations are stored in YAML files and embedded into the Go binary at compile time.

## Supported Languages

| Code | Language |
| ---- | -------- |
| `en` | English |
| `fr` | French |
| `es` | Spanish |
| `it` | Italian |
| `pt` | Portuguese |

## File Structure

```text
i18n/
└── locales/
    ├── locales.go      # Embed directive for all locales
    ├── en/
    │   └── en.yaml     # English translations
    ├── fr/
    │   └── fr.yaml     # French translations
    ├── es/
    │   └── es.yaml     # Spanish translations
    ├── it/
    │   └── it.yaml     # Italian translations
    └── pt/
        └── pt.yaml     # Portuguese translations
```

## Translation File Format

Each YAML file follows this structure:

```yaml
<lang_code>:
  global:
    # Shared translations (buttons, labels, etc.)
    email: "Email"
    password: "Password"
    buttons:
      go-back: "Back"
      delete: "Delete"
  page:
    # Page-specific translations
    home:
      title: "Home"
      h1: "Welcome"
    signin:
      title: "Sign in"
```

### Key Naming Conventions

- Use dot notation for nested keys: `page.home.title`
- Use lowercase with hyphens: `go-back`, `select-language`
- Group by feature/page: `page.signin.*`, `page.match.*`
- Shared translations go in `global.*`

### Variables in Translations

Use `%{variable}` syntax for dynamic values:

```yaml
en:
  page:
    match:
      title: "Match %{name}"
      unauthorized:
        p: "You must be the owner of the match %{matchName} (%{matchId})."
```

## Usage in Code

### In Templ Templates

Import the i18n package and use `i18n.T()`:

```go
import "github.com/invopop/ctxi18n/i18n"

templ MyPage() {
    <h1>{ i18n.T(ctx, "page.home.title") }</h1>
}
```

### With Variables

```go
// Simple variable
i18n.T(ctx, "page.match.title", i18n.M{"name": matchName})

// Multiple variables
i18n.T(ctx, "page.match.unauthorized.p", i18n.M{
    "matchName": match.Name,
    "matchId":   match.ID,
})
```

### In Go Handlers

```go
import "github.com/invopop/ctxi18n/i18n"

func MyHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    message := i18n.T(ctx, "global.error")
}
```

## Adding a New Translation Key

1. Add the key to **all** language files (`en.yaml`, `fr.yaml`, etc.)
2. Use the same nested structure in each file
3. Run `make build/templ` if used in templates
4. Test with different language settings

### Example: Adding a new button

**en/en.yaml:**

```yaml
en:
  global:
    buttons:
      save: "Save"
```

**fr/fr.yaml:**

```yaml
fr:
  global:
    buttons:
      save: "Sauvegarder"
```

## Adding a New Language

1. Create a new directory: `i18n/locales/<code>/`
2. Create the YAML file: `i18n/locales/<code>/<code>.yaml`
3. Copy structure from `en.yaml` and translate all keys
4. Add embed directive to `locales.go`:

```go
//go:embed en
//go:embed fr
//go:embed <code>  // Add this line
var Content embed.FS
```

1. Add the language to the database `language` table

## How Language Selection Works

1. User's preferred language is stored in session/cookie
2. Middleware in `main.go` loads the locale into context:

```go
ctx, err := ctxi18n.WithLocale(r.Context(), lang)
```

1. All `i18n.T()` calls use this context to return the correct translation

## Best Practices

- Always add translations to ALL language files at once
- Keep translation keys descriptive but concise
- Group related translations under common prefixes
- Use variables for dynamic content, never concatenate strings
- Test UI in all languages to check for layout issues (text length varies)
