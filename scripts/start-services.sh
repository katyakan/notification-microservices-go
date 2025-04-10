#!/bin/bash

# Start all services locally for development

echo "🚀 Starting Kafka Notification System services..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  Warning: .env file not found. Creating example .env file..."
    cat > .env << EOF
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
KAFKA_BROKERS=localhost:9092
ENVIRONMENT=development
EOF
    echo "📝 Please edit .env file with your Telegram bot token"
fi

# Function to start a service in a new terminal
start_service() {
    local service_name=$1
    local port=$2
    local cmd=$3
    
    echo "Starting $service_name on port $port..."
    
    # For different operating systems
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        osascript -e "tell application \"Terminal\" to do script \"cd $(pwd) && PORT=$port $cmd\""
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if command -v gnome-terminal &> /dev/null; then
            gnome-terminal --tab --title="$service_name" -- bash -c "cd $(pwd) && PORT=$port $cmd; exec bash"
        elif command -v xterm &> /dev/null; then
            xterm -T "$service_name" -e "cd $(pwd) && PORT=$port $cmd; bash" &
        else
            echo "Please install gnome-terminal or xterm to run services in separate terminals"
            echo "Or run manually: PORT=$port $cmd"
        fi
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        # Windows (Git Bash/Cygwin)
        start powershell -NoExit -Command "cd '$(pwd)'; \$env:PORT='$port'; $cmd"
    else
        echo "Unsupported OS. Please run manually: PORT=$port $cmd"
    fi
}

# Start services
start_service "Producer Service" "3000" "go run ./cmd/producer-service"
sleep 2
start_service "Consumer Service" "3001" "go run ./cmd/consumer-service"
sleep 2
start_service "Notification Service" "3002" "go run ./cmd/notification-service"

echo ""
echo "✅ All services started!"
echo ""
echo "🌐 Producer Service (with Swagger): http://localhost:3000/api/index.html"
echo "⚙️  Consumer Service: http://localhost:3001/health"
echo "📱 Notification Service: http://localhost:3002/health"
echo ""
echo "📋 To test the system, send a POST request to:"
echo "   curl -X POST http://localhost:3000/messages \\"
echo "   -H \"Content-Type: application/json\" \\"
echo "   -d '{\"type\": \"notification\", \"payload\": {\"chatId\": YOUR_CHAT_ID, \"text\": \"Hello from Go!\"}}'"
echo ""
echo "🛑 To stop all services, close the terminal windows or press Ctrl+C in each"
