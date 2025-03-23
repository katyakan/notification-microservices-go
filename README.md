# Kafka Notification System

A microservice-based system for sending notifications to Telegram using Kafka. The system consists of three services:
- Producer Service - accepts HTTP requests and sends messages to Kafka
- Consumer Service - reads messages from Kafka (demonstration service)
- Notification Service - reads messages from Kafka and sends them to Telegram

## Technologies

- Node.js
- NestJS
- Apache Kafka
- Docker
- Telegram Bot API

## Prerequisites

- Docker and Docker Compose
- Node.js 18+
- Telegram Bot Token (get it from @BotFather)

## Installation and Setup

1. Clone the repository:
```bash
git clone git@github.com:merkushov/kafka-notification-system.git
cd kafka-notification-system
```

2. Create `.env` file in the root directory:
```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
```

3. Start the system using Docker Compose:
```bash
docker-compose up -d
```

## Usage

### Sending a Notification

Send a POST request to the Producer Service:

```bash
curl -X POST http://localhost:3000/messages \
-H "Content-Type: application/json" \
-d '{
  "type": "notification",
  "payload": {
    "chatId": YOUR_CHAT_ID,
    "text": "Hello from Producer Service!"
  }
}'
```

### Getting Your Chat ID

1. Find your bot in Telegram
2. Send the `/start` command
3. The bot will reply with your Chat ID

## Architecture

```
┌────────────────┐     ┌─────────┐     ┌────────────────┐     ┌─────────────┐
│ Producer       │     │         │     │ Notification   │     │   Telegram  │
│ Service (3000) │ ──► │  Kafka  │ ──► │ Service       │ ──► │   Bot API   │
└────────────────┘     │         │     └────────────────┘     └─────────────┘
                       │         │     ┌────────────────┐
                       └─────────┘ ──► │ Consumer       │
                                      │ Service (3001)  │
                                      └────────────────┘
```

### System Components

- **Producer Service (port 3000)**
  - Accepts HTTP requests
  - Validates incoming messages
  - Sends messages to Kafka

- **Kafka**
  - Provides asynchronous communication between services
  - Stores messages in the "notifications" topic
  - Supports fault tolerance through the "dead-letter" topic

- **Notification Service**
  - Reads messages from Kafka
  - Sends notifications to Telegram
  - Handles delivery errors

- **Consumer Service (port 3001)**
  - Demonstration service
  - Reads and logs messages from Kafka

## API Documentation

Swagger UI is available at: `http://localhost:3000/api`

## Development

### Project Structure

```
apps/
├── producer-service/     # Message intake service
├── consumer-service/     # Demo consumer
└── notification-service/ # Telegram delivery service
libs/
└── shared/              # Shared interfaces and utilities
```

### Running in Development Mode

```bash
# Install dependencies
npm install

# Start a service
npm run start:dev
```

### Testing

```bash
# unit tests
npm run test

# e2e tests
npm run test:e2e

# test coverage
npm run test:cov
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| TELEGRAM_BOT_TOKEN | Your Telegram bot token | - |
| KAFKA_BROKERS | Kafka broker addresses | localhost:9092 |
| port | Service port | 3000 |

## Docker Services

| Service | Port | Description |
|---------|------|-------------|
| producer-service | 3000 | HTTP API endpoint |
| consumer-service | 3001 | Demo consumer service |
| notification-service | - | Telegram notification sender |
| kafka | 9092 | Message broker |
| zookeeper | 2181 | Kafka dependency |

## Monitoring

- Service logs are available through Docker:
```bash
docker-compose logs -f [service-name]
```

## Local Development Setup

1. Start Kafka infrastructure:
```bash
docker-compose up -d zookeeper kafka
```

2. Run services locally:
```bash
# Terminal 1
npm run start:dev producer-service

# Terminal 2
npm run start:dev consumer-service

# Terminal 3
npm run start:dev notification-service
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[UNLICENSED]