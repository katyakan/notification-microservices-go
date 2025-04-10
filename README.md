# Kafka Notification System (Go)

Микросервисная система для отправки уведомлений в Telegram с использованием Apache Kafka, переписанная на Go. Состоит из трёх сервисов:

- **Producer Service** — принимает HTTP-запросы и отправляет сообщения в Kafka
- **Consumer Service** — получает и логирует сообщения из Kafka (демо-сервис)
- **Notification Service** — получает сообщения из Kafka и отправляет уведомления в Telegram

## Технологии

- Go 1.21+
- Gin Web Framework
- Apache Kafka (segmentio/kafka-go)
- Docker & Docker Compose
- Telegram Bot API (go-telegram-bot-api)
- Swagger/OpenAPI (swaggo)
- Zap Logger
- Viper (конфигурация)

## Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Producer       │    │     Apache      │    │   Notification  │
│  Service        │───▶│     Kafka       │───▶│   Service       │
│  (HTTP API)     │    │                 │    │  (Telegram Bot) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Consumer      │
                       │   Service       │
                       │  (Demo Logger)  │
                       └─────────────────┘
```

## Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)
- Telegram Bot Token (получить у @BotFather)

## Установка и настройка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd kafka-notification-system
```

2. Создайте `.env` файл в корневой директории:
```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
KAFKA_BROKERS=localhost:9092
ENVIRONMENT=development
```

3. Запустите систему с помощью Docker Compose:
```bash
docker-compose up -d
```

Или используйте Makefile:
```bash
make docker-up
```

## Локальная разработка

### Установка зависимостей
```bash
go mod download
```

### Запуск только Kafka и Zookeeper
```bash
make dev-kafka
```

### Запуск сервисов локально

**Вариант 1: Используя скрипты**
```bash
# Linux/macOS
./scripts/start-services.sh

# Windows
scripts\start-services.bat
```

**Вариант 2: Используя Makefile**
```bash
# В разных терминалах
make run-producer
make run-consumer
make run-notification
```

**Вариант 3: Вручную**
```bash
# Терминал 1 - Producer Service
PORT=3000 go run ./cmd/producer-service

# Терминал 2 - Consumer Service
PORT=3001 go run ./cmd/consumer-service

# Терминал 3 - Notification Service
PORT=3002 go run ./cmd/notification-service
```

## Использование

### API Документация

После запуска Producer Service, Swagger документация доступна по адресу:
- http://localhost:3000/api/index.html

### Отправка уведомления

Отправьте POST запрос в Producer Service:

```bash
curl -X POST http://localhost:3000/messages \
-H "Content-Type: application/json" \
-d '{
  "type": "notification",
  "payload": {
    "chatId": YOUR_TELEGRAM_CHAT_ID,
    "text": "Привет от Go Producer Service!"
  }
}'
```

### Health Check

Проверьте статус сервисов:
```bash
curl http://localhost:3000/health  # Producer
curl http://localhost:3001/health  # Consumer
curl http://localhost:3002/health  # Notification
```

## Получение Chat ID

1. Найдите бота [@userinfobot](https://t.me/userinfobot) в Telegram
2. Напишите ему `/start`
3. Скопируйте ваш Chat ID

## Команды Makefile

```bash
make help              # Показать все доступные команды
make build             # Собрать все сервисы
make test              # Запустить тесты
make docker-up         # Запустить с Docker Compose
make docker-down       # Остановить Docker контейнеры
make docker-logs       # Показать логи Docker
make clean             # Очистить артефакты сборки
```

## Структура проекта

```
.
├── cmd/                          # Основные приложения
│   ├── producer-service/         # Producer Service
│   ├── consumer-service/         # Consumer Service
│   └── notification-service/     # Notification Service
├── pkg/                          # Общие пакеты
│   ├── shared/                   # Общие типы и утилиты
│   ├── config/                   # Конфигурация
│   └── logger/                   # Логирование
├── scripts/                      # Скрипты запуска
├── docker-compose.yml            # Docker Compose конфигурация
├── Makefile                      # Команды сборки и запуска
└── go.mod                        # Go модуль
```

## Переменные окружения

| Переменная           | Описание                         | Значение по умолчанию |
|----------------------|----------------------------------|------------------------|
| `TELEGRAM_BOT_TOKEN` | Токен Telegram-бота              | —                      |
| `KAFKA_BROKERS`      | Адреса Kafka-брокеров            | localhost:9092         |
| `PORT`               | Порт сервиса                     | 3000/3001/3002         |
| `ENVIRONMENT`        | Окружение (development/production)| development           |

## Тестирование

```bash
# Запустить все тесты
make test

# Запустить тесты с покрытием
make test-coverage

# Запустить только unit тесты (без интеграционных)
go test -short ./...
```

## Разработка

### Форматирование кода
```bash
make fmt
```

### Линтинг
```bash
make lint
```

### Установка инструментов разработки
```bash
make install-tools
```

## Docker

### Сборка образов
```bash
make docker-build
```

### Запуск отдельных сервисов
```bash
docker-compose up producer-service
docker-compose up consumer-service
docker-compose up notification-service
```

## Мониторинг и логи

Все сервисы используют структурированное логирование с помощью Zap. Логи включают:
- Уровень лога (INFO, ERROR, WARN)
- Временные метки
- Контекстную информацию (ID сообщений, типы, ошибки)

Просмотр логов:
```bash
make docker-logs
```

## Troubleshooting

### Проблемы с Kafka
1. Убедитесь, что Kafka и Zookeeper запущены
2. Проверьте, что топики созданы:
```bash
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Проблемы с Telegram
1. Проверьте правильность токена бота
2. Убедитесь, что бот не заблокирован
3. Проверьте правильность Chat ID

### Проблемы с портами
Убедитесь, что порты 3000, 3001, 3002 свободны:
```bash
netstat -tulpn | grep :300
```

## Лицензия

MIT License
