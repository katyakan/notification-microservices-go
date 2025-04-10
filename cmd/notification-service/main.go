package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kafka-notification-system/cmd/notification-service/internal/handler"
	"kafka-notification-system/cmd/notification-service/internal/service"
	"kafka-notification-system/pkg/config"
	"kafka-notification-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Notification Service API
// @version 1.0
// @description The Notification Service API for Kafka Notification System
// @host localhost:3002
// @BasePath /
func main() {
	// Загружаем конфигурацию
	appConfig := config.LoadAppConfig()
	if appConfig.Port == "3000" {
		appConfig.Port = "3002" // Устанавливаем порт по умолчанию для notification
	}
	kafkaConfig := config.LoadKafkaConfig("notification-service", "telegram-notification-group")

	// Инициализируем логгер
	logger.InitLogger(appConfig.Environment)
	log := logger.GetLogger()

	// Создаем Telegram сервис
	telegramService, err := service.NewTelegramService(appConfig.TelegramBotToken)
	if err != nil {
		log.Fatal("Failed to create Telegram service", zap.Error(err))
	}

	// Создаем Kafka сервис
	kafkaService := service.NewKafkaService(kafkaConfig, telegramService)
	defer func() {
		if err := kafkaService.Close(); err != nil {
			log.Error("Failed to close Kafka service", zap.Error(err))
		}
	}()

	// Создаем обработчики
	notificationHandler := handler.NewNotificationHandler(log)

	// Настраиваем Gin
	if appConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware для логирования
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())

	// Настраиваем маршруты
	v1 := router.Group("/")
	{
		v1.GET("/health", notificationHandler.Health)
	}

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: router,
	}

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем Kafka consumer в горутине
	go func() {
		log.Info("Starting Notification Service Kafka consumer",
			zap.String("topic", kafkaConfig.NotificationsTopic),
			zap.String("groupId", kafkaConfig.GroupID))

		if err := kafkaService.StartConsuming(ctx); err != nil && err != context.Canceled {
			log.Error("Kafka consumer error", zap.Error(err))
		}
	}()

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Info("Starting Notification Service",
			zap.String("port", appConfig.Port),
			zap.Strings("kafka_brokers", kafkaConfig.Brokers))

		fmt.Printf("📱 Notification Service running on port %s\n", appConfig.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Notification Service...")

	// Отменяем контекст для остановки Kafka consumer
	cancel()

	// Graceful shutdown HTTP сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Notification Service stopped")
}
