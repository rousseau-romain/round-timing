# RoundTiming readme

```bash
go install github.com/cespare/reflex@latest
```

## Run

Runs the app and looks for changes.

```bash
make templ
make tailwind
air
```

## Migrate

Create migration file

```bash
migrate create -ext sql -dir database/migration/ -seq {init_message}
```

Run migration up

```bash
make migration_up
```

Run rollback

```bash
make migration_down
```

Run migration fix

```bash
make migration_fix VERSION={VERSION_NUMBER}
```
