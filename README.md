# Livwell

Livwell by Tipen is an e-wallet backend that lets users top up, transfer, pay with QRIS, and track their expenses. It is split into two services, core-service and payment-service, because payment handles much more traffic than the other features and needs to scale on its own. Built with Go, Gin, PostgreSQL, Redis, and RabbitMQ.

## Features

- **Top up:** add balance to a wallet.
- **Transfer:** send money from one wallet to another.
- **Pay:** pay with QRIS, processed asynchronously through payment-service.
- **Expense tracking:** view spending history.

## Architecture

```mermaid
flowchart LR
    Client([Client]) -->|1. QRIS payment request + JWT| Core[core-service]
    Core -->|2. verify session| Redis[(Redis)]
    Core -->|3. publish payment message| MQ{{RabbitMQ}}
    MQ -->|4. consume| Pay[payment-service]
    Pay -->|5. call| GW[Payment gateway - simulated]
    Pay -->|6. publish result| MQ
    MQ -->|7. consume result| Core
    Core -->|8. record transaction| DB[(PostgreSQL)]
```

### How a payment works

1. **core-service** receives the QRIS payment request and publishes a payment message to RabbitMQ.
2. **payment-service** consumes the message and sends the payment to the payment gateway (simulated in this project).
3. **payment-service** publishes the payment result back to RabbitMQ.
4. **core-service** consumes the result and records the transaction in PostgreSQL.

### Design decisions

- **Payment is a separate service.** Payment handles much more traffic than top-up and transfer, so it can be scaled on its own.
- **Payment-service only talks to the outside world.** It handles communication with external payment providers (simulated here), so the hardest and least predictable part of the system is isolated from the core logic.
- **Communication goes through RabbitMQ.** A slow payment gateway never blocks core-service.
- **Authentication stays in core-service.** JWT sessions are verified against Redis at the entry point, so payment-service only receives requests that were already verified.

## Getting Started

**Requirements:** Docker and Docker Compose

```bash
git clone <your-repo-url>
cd <your-repo-folder>
docker compose up --build
```

This starts core-service, payment-service, PostgreSQL, Redis, and RabbitMQ together.
