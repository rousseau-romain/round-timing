# Templ Components

## Overview

RoundTiming uses [Templ](https://templ.guide/) for type-safe HTML templating. Templ compiles `.templ` files into Go code, providing compile-time safety and excellent IDE support.

## File Locations

```text
shared/
└── components/          # Reusable UI components
    ├── layout.templ     # Main page layout
    ├── footer.templ     # Site footer
    ├── button-copy.templ
    ├── button-enable-disaled.templ
    ├── svg-heart.templ
    └── svg-user.templ

views/
└── page/               # Full page templates
    ├── home.templ
    ├── signin.templ
    ├── signup.templ
    ├── profile.templ
    ├── match.templ
    ├── match-list.templ
    ├── start-match.templ
    ├── spectate-match.templ
    ├── 403.templ
    ├── 404.templ
    ├── cgu.templ
    ├── privacy.templ
    └── admin/
        └── user-list.templ
```

## Templ Syntax Basics

### Package Declaration

```go
package components

import (
    "github.com/invopop/ctxi18n/i18n"
    "github.com/rousseau-romain/round-timing/model"
)
```

### Defining a Component

```go
templ Button(text string, class string) {
    <button class={ class }>{ text }</button>
}
```

### Using Components

```go
templ MyPage() {
    @Button("Click me", "btn")
}
```

### Children (Slots)

```go
templ Card() {
    <div class="card">
        { children... }
    </div>
}

// Usage
@Card() {
    <p>Content inside the card</p>
}
```

### Conditionals

```go
templ UserStatus(isLoggedIn bool) {
    if isLoggedIn {
        <span>Welcome back!</span>
    } else {
        <a href="/signin">Sign in</a>
    }
}
```

### Loops

```go
templ NavList(items []NavItem) {
    <ul>
        for _, item := range items {
            <li>
                <a href={ templ.SafeURL(item.Url) }>{ item.Name }</a>
            </li>
        }
    </ul>
}
```

### Dynamic Classes

```go
templ Alert(alertType string) {
    <div class={
        "alert",
        templ.KV("alert-success", alertType == "success"),
        templ.KV("alert-error", alertType == "error"),
    }>
        { children... }
    </div>
}
```

## Layout Component

The main layout wraps all pages:

```go
// shared/components/layout.templ
templ Layout(title string, popinMessages PopinMessages, user model.User, navItems []NavItem, languages []model.Language, pageSlug string) {
    <html>
        <head>
            <title>Round Timing - { title }</title>
            // CSS, favicon...
        </head>
        <body>
            // HTMX scripts
            @Nav(user, navItems, languages, pageSlug)
            <div class="min-h-screen" id="content">
                { children... }
            </div>
            @Footer()
        </body>
    </html>
}
```

### Usage in Pages

```go
// views/page/home.templ
templ HomePage(user model.User, popinMessage components.PopinMessages, navItems []components.NavItem, languages []model.Language, pageSlug string) {
    @components.Layout(i18n.T(ctx, "page.home.title"), popinMessage, user, navItems, languages, pageSlug) {
        <div class="container mx-auto p-4">
            // Page content here
        </div>
    }
}
```

## Common Patterns

### Translation Integration

```go
import "github.com/invopop/ctxi18n/i18n"

templ MyComponent() {
    <h1>{ i18n.T(ctx, "page.home.title") }</h1>
    <p>{ i18n.T(ctx, "page.home.welcome", i18n.M{"name": userName}) }</p>
}
```

### Raw HTML (Use Sparingly)

```go
templ RichContent() {
    @templ.Raw(i18n.T(ctx, "page.home.discover", i18n.M{"name": "<strong>Round Timing</strong>"}))
}
```

### HTMX Attributes

```go
templ DeleteButton(id int) {
    <button
        hx-delete={ fmt.Sprintf("/match/%d", id) }
        hx-target="#match-list"
        hx-swap="outerHTML"
    >
        Delete
    </button>
}
```

### Safe URLs

```go
templ Link(url string) {
    <a href={ templ.SafeURL(url) }>Link</a>
}
```

## Generating Templ Files

After editing `.templ` files, generate the Go code:

```bash
# Single generation
templ generate

# Or using make
make build/templ

# Watch mode (auto-regenerate)
templ generate -watch
```

## Best Practices

1. **Keep components small** - Each component should do one thing
2. **Use typed parameters** - Leverage Go's type system
3. **Prefer composition** - Build complex UIs from simple components
4. **Never edit `*_templ.go`** - These are generated files
5. **Use `templ.KV` for conditional classes** - Cleaner than string concatenation
6. **Pass context for i18n** - Always use `ctx` for translations
