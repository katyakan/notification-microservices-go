package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kafka-notification-system/cmd/consumer-service/internal/handler"
	"kafka-notification-system/cmd/consumer-service/internal/service"
	"kafka-notification-system/pkg/config"
	"kafka-notification-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Consumer Service API
// @version 1.0
// @description The Consumer Service API for Kafka Notification System
// @host localhost:3001
// @BasePath /
func main() {
	// Загружаем конфигурацию
	appConfig := config.LoadAppConfig()
	if appConfig.Port == "3000" {
		appConfig.Port = "3001" // Устанавливаем порт по умолчанию для consumer
	}
	kafkaConfig := config.LoadKafkaConfig("consumer-service", "notification-group")


	logger.InitLogger(appConfig.Environment)
	log := logger.GetLogger()


	kafkaService := service.NewKafkaService(kafkaConfig)
	defer func() {
		if err := kafkaService.Close(); err != nil {
			log.Error("Failed to close Kafka service", zap.Error(err))
		}
	}()


	consumerHandler := handler.NewConsumerHandler(log)

	
	if appConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

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
		v1.GET("/health", consumerHandler.Health)
	}


	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: router,
	}

	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	go func() {
		log.Info("Starting Kafka consumer",
			zap.String("topic", kafkaConfig.NotificationsTopic),
			zap.String("groupId", kafkaConfig.GroupID))

		if err := kafkaService.StartConsuming(ctx); err != nil && err != context.Canceled {
			log.Error("Kafka consumer error", zap.Error(err))
		}
	}()


	go func() {
		log.Info("Starting Consumer Service",
			zap.String("port", appConfig.Port),
			zap.Strings("kafka_brokers", kafkaConfig.Brokers))

		fmt.Printf("⚙️ Consumer Service running on port %s\n", appConfig.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Consumer Service...")

	
	cancel()


	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Consumer Service stopped")
}
