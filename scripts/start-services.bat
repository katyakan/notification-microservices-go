@echo off
echo 🚀 Starting Kafka Notification System services...

REM Check if .env file exists
if not exist .env (
    echo ⚠️  Warning: .env file not found. Creating example .env file...
    (
        echo TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
        echo KAFKA_BROKERS=localhost:9092
        echo ENVIRONMENT=development
    ) > .env
    echo 📝 Please edit .env file with your Telegram bot token
)

echo Starting Producer Service on port 3000...
start "Producer Service" powershell -NoExit -Command "$env:PORT='3000'; go run ./cmd/producer-service"

timeout /t 2 /nobreak >nul

echo Starting Consumer Service on port 3001...
start "Consumer Service" powershell -NoExit -Command "$env:PORT='3001'; go run ./cmd/consumer-service"

timeout /t 2 /nobreak >nul

echo Starting Notification Service on port 3002...
start "Notification Service" powershell -NoExit -Command "$env:PORT='3002'; go run ./cmd/notification-service"

echo.
echo ✅ All services started!
echo.
echo 🌐 Producer Service (with Swagger): http://localhost:3000/api/index.html
echo ⚙️  Consumer Service: http://localhost:3001/health
echo 📱 Notification Service: http://localhost:3002/health
echo.
echo 📋 To test the system, send a POST request to:
echo    curl -X POST http://localhost:3000/messages ^
echo    -H "Content-Type: application/json" ^
echo    -d "{\"type\": \"notification\", \"payload\": {\"chatId\": YOUR_CHAT_ID, \"text\": \"Hello from Go!\"}}"
echo.
echo 🛑 To stop all services, close the PowerShell windows
pause
