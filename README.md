# RoundTiming readme

You nedd to have:

- brew
- docker v27.3.1
- docker-compose v2.30.3
- go v1.22.9
- node v23.2.0
- npm 10.9.0

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
make migration_create {init_message}
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
