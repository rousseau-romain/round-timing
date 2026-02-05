# Monitoring

Loki + Promtail + Grafana stack for log aggregation and visualization.

## Start / Stop

```bash
make monitoring_start   # Start Loki, Promtail, Grafana
make monitoring_stop    # Stop monitoring stack
```

## Run the app with log collection

`make live` automatically writes logs to `./logs/app.log` so Promtail ships them to Loki.

```bash
make live
```

## Grafana

Open http://localhost:3000 → **Explore** → select **Loki** datasource.

### LogQL query examples

Show all app logs:

```
{job="round-timing"} | json
```

Filter by message:

```
{job="round-timing"} | json | msg="match created"
{job="round-timing"} | json | msg="user created"
{job="round-timing"} | json | msg=~".*created.*"
```

Filter by level:

```
{job="round-timing"} | json | level="INFO"
{job="round-timing"} | json | level="ERROR"
```

Filter by field value:

```
{job="round-timing"} | json | id_user="1"
```

Combine multiple filters:

```
{job="round-timing"} | json | msg="match created" | id_user="1"
```

Show only action logs (added in handlers):

```
{job="round-timing"} | json | msg=~".*(created|deleted|started|reset|toggled|changed|logged|used|removed).*"
```

Show logs with a specific key:

```
{job="round-timing"} | json | source != ""
```

Show MariaDB container logs:

```
{job="docker"}
```

## Architecture

```
App (slog JSON) → ./logs/app.log → Promtail → Loki → Grafana
MariaDB container → Docker logs   → Promtail → Loki → Grafana
```

## Ports (local)

| Service  | Port | URL                    |
|----------|------|------------------------|
| Grafana  | 3000 | http://localhost:3000   |
| Loki     | 3100 | http://localhost:3100   |
| Promtail | 9080 | http://localhost:9080   |

## Deploy on Coolify

The monitoring stack is deployed as a single **Docker Compose** resource. All services share a network and can communicate by name.

### 1. Create the resource

**New Resource** → **Docker Compose** → select your repo with:

- **Base Directory**: `/monitoring`
- **Docker Compose File**: `docker-compose.monitoring.yml`
- **Branch**: `feat/add-graphana-logs` (or `master` once merged)

### 2. Domain

In the Grafana service settings, set the domain with the container port:

```
https://your-grafana-domain.com:3000
```

The `:3000` tells Coolify's proxy which container port to route to.

### 3. Environment variables

```
GF_ADMIN_PASSWORD=<your-password>
```

### 4. Volumes

Coolify handles persistent volumes for Loki (`/loki`) and Grafana (`/var/lib/grafana`) automatically.

Add these host mounts for Promtail in Coolify's **Storages** tab:

```
/var/lib/docker/containers:/var/lib/docker/containers:ro
/var/run/docker.sock:/var/run/docker.sock:ro
```

### 5. Watch Paths

Set watch path to only redeploy when monitoring config changes:

```
monitoring/**
```

### Configuration files

Configs are baked into Docker images via Dockerfiles (not mounted as volumes):

| Service  | Dockerfile               | Config                                              |
|----------|--------------------------|-----------------------------------------------------|
| Loki     | `loki/Dockerfile`        | `loki/loki-config.yml`                              |
| Promtail | `promtail/Dockerfile`    | `promtail/promtail-prod-config.yml`                 |
| Grafana  | `grafana/Dockerfile`     | `grafana/provisioning/datasources/datasources-prod.yml` |

To change a config, edit the file, push, and redeploy.

### Logging pattern

Handlers use a local `logger` variable to avoid mutating the shared `h.Slog`:

```go
func (h *Handler) HandleCreateMatch(w http.ResponseWriter, r *http.Request) {
    user, _ := auth.UserFromRequest(r)
    logger := h.Slog.With("userId", user.Id)  // local, not h.Slog =

    logger.Info("match created", "matchId", matchId)
    logger.Error(err.Error())
}
```

Using `h.Slog = h.Slog.With(...)` would append `userId` on every request, causing duplicate keys to accumulate.

### Production architecture

```
App container → Docker logs → Promtail → Loki → Grafana
```
