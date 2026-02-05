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

Each monitoring service is deployed as a separate Coolify resource using its own Dockerfile.

### 1. Create the services

For each service, create a **New Resource** → **Public/Private Repository** with:

| Service  | Base Directory         | Dockerfile   | Port |
|----------|------------------------|-------------|------|
| Loki     | `/monitoring/loki`     | `Dockerfile` | 3100 |
| Promtail | `/monitoring/promtail` | `Dockerfile` | 9080 |
| Grafana  | `/monitoring/grafana`  | `Dockerfile` | 3000 |

### 2. Environment variables

**Promtail:**

```
LOKI_URL=http://loki:3100
```

**Grafana:**

```
LOKI_URL=http://loki:3100
GF_SECURITY_ADMIN_PASSWORD=<your-password>
```

### 3. Volumes

| Service  | Mount                                                  |
|----------|--------------------------------------------------------|
| Loki     | Persistent volume → `/loki`                            |
| Promtail | `/var/lib/docker/containers:/var/lib/docker/containers:ro` |
| Promtail | `/var/run/docker.sock:/var/run/docker.sock:ro`         |
| Grafana  | Persistent volume → `/var/lib/grafana`                 |

### 4. Networking

Enable **Connect to Predefined Networks** in **Settings** for all three services and for the main app. This puts them on the same Docker network so they can communicate by service name.

### 5. Watch Paths

Set watch paths so monitoring services only redeploy when their config changes:

| Service  | Watch Path               |
|----------|--------------------------|
| Loki     | `monitoring/loki/**`     |
| Promtail | `monitoring/promtail/**` |
| Grafana  | `monitoring/grafana/**`  |

### Production architecture

```
App container → Docker logs → Promtail → Loki → Grafana
```
