# Deployment

## Environments

| Environment | Branch | URL |
| ------------- | -------- | ----- |
| Production | `master` | TBD |
| Staging | `staging` | <https://round-timing-staging.web-rows.ovh/> |
| Development | local | <http://127.0.0.1:7331> |

## Staging Deployment

Staging deployments are automatic when pushing to the `staging` branch:

```bash
git checkout staging
git merge feature-branch
git push origin staging
```

The staging server will automatically pull and deploy the changes.

## Production Deployment

Production deployments are triggered by pushing to `master`:

```bash
git checkout master
git merge staging
git push origin master
```

## Docker Build

The project includes a `Dockerfile` for containerized deployment:

```bash
docker build -t round-timing .
docker run -p 7331:7331 --env-file .env.prod round-timing
```

## Environment Files

| File | Purpose |
| ------ | --------- |
| `.env` | Local development |
| `.env.staging` | Staging environment |
| `.env.prod` | Production environment |
| `.env.template` | Template for new environments |

## Pre-Deployment Checklist

1. [ ] All tests pass (manual testing)
2. [ ] Templ files generated (`make build/templ`)
3. [ ] CSS built (`make build/tailwind`)
4. [ ] Migrations are encrypted (`make db/encrypt`)
5. [ ] No sensitive data in commits
6. [ ] Environment variables configured on server

## Rollback

To rollback a deployment:

1. Identify the last working commit
2. Reset the branch: `git reset --hard <commit>`
3. Force push: `git push --force origin <branch>`

For database rollbacks:

```bash
make migration_down  # Rollback one migration
```
