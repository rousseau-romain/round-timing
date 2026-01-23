# Claude Code Instructions for RoundTiming

## Project Overview

RoundTiming is a Go web application for tracking game rounds/timing (likely for esports/gaming). It uses server-side rendering with Templ templates and TailwindCSS.

## Tech Stack

- **Go 1.25** with Gorilla Mux router
- **Templ** for HTML templating (`.templ` files compile to `_templ.go`)
- **TailwindCSS** for styling
- **MySQL** database via Docker
- **OAuth** authentication (Discord, Google) via Goth

## Development Commands

```bash
make live          # Start dev server with hot reload (localhost:7331)
make db_start      # Start database container
make db_stop       # Stop database container
make migration_up  # Run database migrations
make migration_down # Rollback last migration
make build/templ   # Generate templ files
make build/tailwind # Build CSS
```

## Project Structure

- `handlers/` - HTTP route handlers
- `model/` - Database models and queries
- `views/page/` - Page templates (`.templ` files)
- `views/components/` - Reusable UI components organized by type:
  - `layout/` - Layout, Footer, Nav, PopinMessages
  - `ui/` - Buttons, AvatarToggle, ErrorPageContent
  - `icons/` - SVG icon components (Heart, User)
  - `forms/` - Form input components
- `service/` - Business logic
- `config/` - Configuration loading
- `i18n/` - Internationalization

## Code Patterns

### Templ Templates

- Template files use `.templ` extension
- Run `templ generate` after modifying `.templ` files (or use `make live`)
- Generated files are `*_templ.go` - do not edit these directly
- Import components by package:
  ```go
  import "github.com/rousseau-romain/round-timing/views/components/layout"
  import "github.com/rousseau-romain/round-timing/views/components/ui"
  ```

### UI Components (`views/components/ui/`)

**Buttons** (`button.templ`):
```go
// Basic button
@ui.Button("primary", "md") { Click me }

// Button with HTMX attributes
@ui.ButtonAction("danger", "sm", templ.Attributes{
    "hx-delete": "/item/1",
}) { Delete }

// Link styled as button
@ui.ButtonLink("indigo", "lg", "/path") { Go }

// Link with custom attributes
@ui.ButtonLinkAction("black", "lg", "/path", templ.Attributes{
    "class": "w-full",
}) { Submit }
```

Variants: `primary` (blue), `success` (green), `danger` (red), `outline`, `indigo`, `black`
Sizes: `sm`, `md`, `lg`

**Badges** (`badge.templ`):
```go
// Generic badge with variant
@ui.Badge("red", templ.Attributes{"class": "absolute -right-1 -bottom-1"}) { 3 }

// Recovery badge (auto-colors based on rounds: 1=red, 2=yellow, 3+=cyan)
@ui.BadgeRecovery(mps.RoundBeforeRecovery)
```

Variants: `red`, `yellow`, `cyan`, `green`, `indigo`, `gray`

**Tables** (`table.templ`):
```go
@ui.Table("default") {
    @ui.TableHead("default") {
        <tr>
            @ui.Th() { Name }
            @ui.ThEmpty()
        </tr>
    }
    @ui.TableBody(templ.Attributes{
        "hx-swap": "outerHTML",
        "hx-target": "closest tr",
    }) {
        <tr>
            @ui.TdPrimary() { Item name }
            @ui.TdAction() { @ui.Button(...) }
        </tr>
    }
}
```

Table variants: `default` (full-width dividers), `compact` (bordered auto-width)
Rows: `Tr`, `TrBorder`, `TrColor(color)`
Header cells: `Th`, `ThEmpty`, `ThCompact`
Body cells: `Td`, `TdPrimary`, `TdCenter`, `TdAction`, `TdCompact`

### Form Components (`views/components/forms/`)

**forms.templ**:
```go
// Input with label (grid layout)
@forms.Input("email", "email", "email", "global.email", true)

// Flexible input with custom attributes
@forms.InputAction("text", "name", "Label", templ.Attributes{
    "placeholder": "Enter name",
    "hx-post":     "/update",
})

// Input without label (pass empty string)
@forms.InputAction("text", "search", "", templ.Attributes{...})

// Select dropdown
@forms.Select("country", "country", "form.country", []forms.SelectOption{
    {Value: "fr", Label: "France", Selected: true},
    {Value: "us", Label: "USA"},
}, true)

// Textarea
@forms.Textarea("bio", "bio", "form.bio", 4, false)

// Checkbox
@forms.Checkbox("terms", "terms", "form.accept-terms", false)

// Radio group
@forms.Radio("gender", []forms.SelectOption{...}, "gender", true)
```

### Styling with TailwindCSS

- Configuration: `tailwind.config.js`
- Input CSS: `input.css` (organized with section comments)
- Dynamic classes for team colors use safelist patterns in config
- CSS variables for theming defined in `:root` and `.dark`
- CSS organization:
  - CSS Variables (design tokens)
  - Typography (h1-h3)
  - Form elements (inputs, selects, checkboxes)
  - Links (content, breadcrumbs, footer)
  - Utility classes (tooltip)
  - HTMX integration (swap animations)

### Database

- Models are in `model/` directory
- Uses `huandu/go-sqlbuilder` for query building
- Migrations are in `database/migration/` (encrypted with GPG)

### Authentication

- Session-based auth via `gorilla/sessions`
- OAuth providers: Discord, Google
- JWT for API tokens

## Testing

Run the application manually to test changes since there are no automated tests visible.

## Git Workflow

- `master` - Production branch
- `staging` - Staging branch (auto-deploys to staging URL)
- Push to `staging` to deploy to <https://round-timing-staging.web-rows.ovh/>

## Environment Variables

Copy `.env.template` to `.env` and configure:

- Database connection (DB_*)
- OAuth credentials (DISCORD_*, GOOGLE_*)
- Session secrets (COOKIES_AUTH_*, JWT_SECRET_KEY, SALT_SECRET)
