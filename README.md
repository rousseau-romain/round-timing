# RoundTiming readme

You nedd to have:

- brew
- docker
- docker-compose

## Run

Runs the app and looks for changes.

```bash
make db_start
make install
make live
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
