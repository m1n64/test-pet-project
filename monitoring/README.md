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

<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/dbd384e6-f83d-4cd1-ab1f-759ff8cdeff8" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/102861ae-fa90-44eb-972e-3565fa618544" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/5a60330a-3c81-47e5-ba7f-4906f41e496f" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/02a7304e-b513-41dd-a79f-f6bad2abf0d5" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/5a4cbe70-aa0f-4614-8442-5a7ba433c879" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/65ed7e19-77fd-44d1-ad32-60cb7bc90c12" />


