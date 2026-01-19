# Development Setup

## Prerequisites

Install the following tools:

| Tool | Version | Installation |
| ---- | ------- | ------------ |
| Homebrew | latest | [brew.sh](https://brew.sh) |
| Docker | v27.3.1+ | `brew install docker` |
| Docker Compose | v2.30.3+ | Included with Docker |
| Go | v1.25.6 | `brew install go` |
| Node.js | v23.2.0 | `nvm install 23.2.0` (recommended) |
| npm | v10.9.0 | Included with Node.js |

### Recommended: Version Managers

- **Node**: Use [nvm](https://github.com/nvm-sh/nvm) for Node version management

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/rousseau-romain/round-timing.git
cd round-timing
```

### 2. Environment Configuration

```bash
cp .env.template .env
```

Edit `.env` with your configuration:

```text
# Database
DB_HOST=localhost
DB_PORT=3306
DB_DRIVER=mysql
DB_NAME=round_timing
DB_USER=user
DB_PASSWORD=password
DB_ROOT_PASSWORD=root_password

# Auth cookies
COOKIES_AUTH_SECRET=<generate-random-string>
COOKIES_AUTH_AGE_IN_SECONDS=86400
COOKIES_AUTH_IS_SECURE=false
COOKIES_AUTH_IS_HTTP_ONLY=true

# OAuth (optional for local dev)
DISCORD_CLIENT_ID=<your-discord-client-id>
DISCORD_CLIENT_SECRET=<your-discord-client-secret>
GOOGLE_CLIENT_ID=<your-google-client-id>
GOOGLE_CLIENT_SECRET=<your-google-client-secret>

# Security
SALT_SECRET=<generate-random-string>
JWT_SECRET_KEY=<generate-random-string>

# Server
PUBLIC_HOST_PORT=http://127.0.0.1:7331
```

### 3. Decrypt Database Migrations

You need the GPG passphrase to decrypt migrations:

```bash
make db/decrypt
```

### 4. Install Dependencies

```bash
make install
```

This installs:

- `golang-migrate` - Database migrations
- `gnupg` - Encryption for migrations
- `air` - Go live reload
- `templ` - Template compiler
- npm packages (TailwindCSS)

### 5. Initialize Database

```bash
make db_init
```

This starts Docker and runs migrations.

### 6. Encrypt Migrations (Important!)

After database setup, re-encrypt the migrations:

```bash
make db/encrypt
```

## Running the Development Server

```bash
make live
```

This starts:

- **Go server** with hot reload (Air) on port 2468
- **Templ watcher** with proxy on port 7331
- **TailwindCSS** watcher

Access the app at: **<http://127.0.0.1:7331>**

## Development Workflow

### Making Changes

1. **Go code**: Save file, Air auto-reloads
2. **Templ files**: Save file, Templ auto-regenerates
3. **CSS/TailwindCSS**: Save file, Tailwind auto-rebuilds

### Database Changes

```bash
# Create a new migration
make migration_create add_new_table

# Run migrations
make migration_up

# Rollback
make migration_down
```

### Building for Production

```bash
make build/templ      # Generate templ files
make build/tailwind   # Build minified CSS
go build -o app       # Compile binary
```

## Common Commands

| Command | Description |
| ------- | ----------- |
| `make live` | Start dev server with hot reload |
| `make db_start` | Start database container |
| `make db_stop` | Stop database container |
| `make migration_up` | Run pending migrations |
| `make migration_down` | Rollback last migration |
| `make build/templ` | Generate templ files |
| `make build/tailwind` | Build CSS |
| `make show_deadcode` | Find unused code |

## VS Code Setup

Add to your `.vscode/settings.json`:

```json
{
  "go.goroot": "$GOROOT"
}
```

Install recommended extensions:

- Go (golang.go)
- Templ (a-h.templ)
- Tailwind CSS IntelliSense

## Troubleshooting

### Port Already in Use

```bash
# Find and kill process on port 2468
lsof -i :2468 | grep LISTEN | awk '{print $2}' | xargs kill
```

### Database Connection Failed

```bash
# Check if Docker is running
docker ps

# Restart database
make db_stop
make db_start
```

### Templ Not Regenerating

```bash
# Manual regenerate
templ generate

# Check for syntax errors in .templ files
```

### Migration Errors

```bash
# Fix dirty migration state
make migration_fix VERSION=<last_good_version>
make migration_down
# Fix the SQL
make migration_up
```
