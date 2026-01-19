# RoundTiming readme

You nedd to have:

- brew
- docker v27.3.1
- docker-compose v2.30.3
- go v1.25.6 (brew install go)
- node v23.2.0 (nvm is better)
- npm 10.9.0

## Run

Create your `.env` based on `env.template` you can run:

```bash
cp .env.template .env
```

Run these commands first(you need the key for the ecryption):

```bash
make db/decrypt
make db_start
make db_init
make db/encrypt
make install
make live
```

Runs the app and looks for changes on `127.0.0.1:7331` for live reload.

## Migrate

### Create migration file

```bash
make migration_create {init_message}
```

### Run migration up

```bash
make migration_up
```

### Run rollback

```bash
make migration_down
```

### Run migration fix

You did a shity make migration up and you need to fix it don't wory do theses steps (for VERSION_NUMBER=18):

```bash
make migration_fix VERSION={VERSION_NUMBER}
make migration_down
```

Now you are in state VERSION_NUMBER - 1.
If your rollback ix the database you can now run:

```bash
make migration_up
```

If not you can debug and fix !

## Deploy

Push to branch staging to deploy on this url https://round-timing-staging.web-rows.ovh/
