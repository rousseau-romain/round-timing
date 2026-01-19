# /dev - Development Server

Start the development server with hot reload.

## Usage

```text
/dev        # Start dev server
/dev stop   # Stop all dev processes
```

## Instructions

When the user runs `/dev`:

1. Check if Docker is running with `docker ps`
2. Start the database if not running: `make db_start`
3. Run `make live` to start the dev server with hot reload
4. The server runs on `http://127.0.0.1:7331` with Templ proxy on port 2468

When the user runs `/dev stop`:

1. Stop any running `air` processes
2. Optionally stop the database with `make db_stop`
