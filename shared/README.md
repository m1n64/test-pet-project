# Shared service

This service contains Message Broker and other shared services.

## Services
- **Message Broker**: RabbitMQ
- **Mail**: MailHog (For development only)

## Run
```bash
cp .env.example .env
```

```bash
make up
```

## Ports
- 5671: RabbitMQ
- 8025: MailHog
- 15671: RabbitMQ Management UI
- 1025: SMTP (MailHog)

## Screenshots
