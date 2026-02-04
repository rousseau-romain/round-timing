# /db - Database Management

Manage the MySQL database container and data.

## Usage

```text
/db start   # Start database container
/db stop    # Stop database container
/db init    # Initialize database (first time setup)
/db reset   # Reset database (recreate from migrations)
```

## Instructions

### /db start

Run `make db_start` to start the Docker MySQL container.

### /db stop

Run `make db_stop` to stop the Docker container.

### /db init

Run `make db_init` which:

1. Starts Docker container with `docker-compose up -d`
2. Waits for MySQL to be ready
3. Runs all migrations with `make migration_up`

### /db reset

1. Stop the database: `make db_stop`
2. Remove the container and volume: `docker-compose down -v`
3. Reinitialize: `make db_init`

Note: Database migrations are encrypted with GPG. Run `make db/decrypt` before working with migrations if needed.
