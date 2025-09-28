# Monitoring service

This service is responsible for monitoring and logs management. It collects logs from various services, analyzes them for errors or performance issues, and provides a dashboard for real-time monitoring.

## Services

**Logs Collector** - Collects logs from different services and stores them in a centralized database:
- MongoDB
- OpenSearch
- Graylog

**Monitoring Dashboard** - Provides a web interface for real-time monitoring and visualization of logs:
- Grafana
- InluxDB
- Telegraf
- docker-socket-proxy - for monitoring Docker containers

## Run
```bash
cp .env.example .env
```

```bash
make up
```

In **Graylog** you can create a UDP input on port `12201` to receive logs in GELF format (`System` -> `Inputs` -> `Select input` -> `GELF UDP`).

## Ports
- 3000 - Grafana
- 9000 - Graylog
- 9200 - OpenSearch
- 8086 - InfluxDB
- 12201 - GELF UDP (Graylog input)
- 12202 - GELF TCP (Graylog input)
- 5555 - Graylog
- 8090 - Telegraf UDP

## Screenshots