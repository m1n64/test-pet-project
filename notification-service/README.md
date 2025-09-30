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

To generate structures for cache:
```bash
make generate # optional
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

<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/c7433f01-6185-47a5-ad0d-1bfd6db68a26" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/1d23e61e-3906-4886-a5a2-3a1abd6c8841" />
<img width="1483" height="952" alt="image" src="https://github.com/user-attachments/assets/b80a9876-f145-4ce9-8afe-ff1e8d3c922b" />
<img width="1484" height="952" alt="image" src="https://github.com/user-attachments/assets/1c37003c-6d98-4e93-8bec-b7cc5b84549e" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/d0ee2063-6c17-4203-b8e1-850e543594c5" />
<img width="1988" height="1047" alt="image" src="https://github.com/user-attachments/assets/7fc98bd0-0386-4b57-a4af-5ce145327d29" />

