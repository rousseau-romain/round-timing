# RoundTiming readme

You nedd to have:

- brew
- docker v27.3.1
- docker-compose v2.30.3
- go v1.22.9 (gvm is better)
- node v23.2.0
- npm 10.9.0 (nvm is better)

## Run

Create your `.env` based on `env.template` you can run:

```bash
cp .env.template .env
```

Run these commands first

```bash
make db_start
make install
make live
```

Runs the app and looks for changes on `127.0.0.1:7331` for live reload.

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
