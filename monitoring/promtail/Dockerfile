FROM grafana/promtail:3.4

COPY promtail-prod-config.yml /etc/promtail/config.yml

CMD ["-config.file=/etc/promtail/config.yml", "-config.expand-env=true"]
