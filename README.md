# RoundTiming readme

You nedd to have:

- golang-migrate gnupg
- docker v27.3.1
- docker-compose v2.30.3
- go v1.25.6 (brew install go)
- node v23.2.0 (nvm is better)
- npm 10.9.0

Optional (for releasing):

- git-cliff (changelog generation)
- goreleaser (GitHub releases)

## Run

Create your `.env` based on `env.template` you can run:

```bash
cp .env.template .env
```

Run these commands first(you need the key for the ecryption):

```bash
make db/decrypt
make install
make db_init
make db_start
make db/encrypt
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

## Release

### Prerequisites

```bash
# Arch Linux
yay -S git-cliff goreleaser-bin
```

### Usage

```bash
make release          # auto-detect bump type from conventional commits, update CHANGELOG.md, commit, tag
make release/major    # force major bump
make release/minor    # force minor bump
make release/patch    # force patch bump
make release/push     # git push origin master --tags
make release/github   # create GitHub release (requires GITHUB_TOKEN)
make changelog        # regenerate full CHANGELOG.md
```

The `release` target analyzes commits since the last tag:
- `BREAKING CHANGE` or `!:` → major bump
- `feat:` → minor bump
- `fix:` → patch bump

## Deploy

Push to branch staging to deploy on this url https://round-timing-staging.web-rows.ovh/

### App (Coolify)

The app is deployed as a **Dockerfile** resource on Coolify.

1. **New Resource** → **Dockerfile** → select the repo
2. **Branch**: `staging` or `master`
3. **Port**: `2468`
4. **Domain**: `https://your-domain.com:2468`
5. **Environment variables**: copy from `.env.template` and fill in values

### Monitoring (Coolify)

See [monitoring/README.md](monitoring/README.md) for Loki + Grafana deployment instructions.
