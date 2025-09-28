# Notification service

Microservice for sending notifications to users.

## Tech stack
- Go 1.25
- Gin
- Postgres
- Redis

## Run
```bash
cp .env.example .env
```

```bash
make up
```

## Ports
- 5878 - JSON-RPC
- 5864 - Debug
- 5432 - Postgres
- 6379 - Redis

## Channels
- Telegram
- Email
- SMS (WIP)
- Push (WIP)

## API

**JSON-RPC** specification:
- http://localhost:5878

## Screenshots