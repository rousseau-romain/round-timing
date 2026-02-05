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

## Ports

| Service  | Port | URL                    |
|----------|------|------------------------|
| Grafana  | 3000 | http://localhost:3000   |
| Loki     | 3100 | http://localhost:3100   |
| Promtail | 9080 | http://localhost:9080   |
