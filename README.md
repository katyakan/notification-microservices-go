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


2. Create `.env` file in the root directory:
```env
TELEGRAM_BOT_TOKEN=<your_telegram_bot_token_here>
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
    "chatId": YOUR_TELEGRAM_CHAT_ID,
    "text": "Hello from Producer Service!"
  }
}'
```
