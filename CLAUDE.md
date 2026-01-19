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
- `shared/components/` - Reusable UI components
- `service/` - Business logic
- `config/` - Configuration loading
- `i18n/` - Internationalization

## Code Patterns

### Templ Templates

- Template files use `.templ` extension
- Run `templ generate` after modifying `.templ` files (or use `make live`)
- Generated files are `*_templ.go` - do not edit these directly

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
