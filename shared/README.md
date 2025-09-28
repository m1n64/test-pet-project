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

<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/c3c6f42c-7c06-4cd9-af60-e45df07799f0" />
<img width="910" height="512" alt="image" src="https://github.com/user-attachments/assets/fe2970e1-0522-4da8-96f3-ac1bfb298dfd" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/588cc1d2-1043-4905-a976-be0be1ede27e" />
<img width="2032" height="1091" alt="image" src="https://github.com/user-attachments/assets/592f7ce9-9ae2-4d45-9cd0-c8411ea6b9f2" />
