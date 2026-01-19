# /migrate - Database Migrations

Manage database schema migrations.

## Usage

```text
/migrate up              # Run all pending migrations
/migrate down            # Rollback last migration
/migrate create <name>   # Create new migration files
/migrate fix <version>   # Force migration version (for fixing issues)
```

## Instructions

### /migrate up

Run `make migration_up` to apply all pending migrations.

### /migrate down

Run `make migration_down` to rollback the last migration.

### /migrate create {name}

Run `make migration_create <name>` to create new migration files.
This creates two files in `database/migration/`:

- `NNNN_<name>.up.sql` - Forward migration
- `NNNN_<name>.down.sql` - Rollback migration

### /migrate fix {version}

If a migration fails partially, use this to fix the state:

1. Run `make migration_fix VERSION=<version>`
2. Run `make migration_down` to rollback
3. Fix the migration SQL
4. Run `make migration_up` to retry

## Important Notes

- Migrations are encrypted with GPG for security
- Run `make db/decrypt` before editing migrations
- Run `make db/encrypt` after editing to secure them
- Never commit decrypted migration files
