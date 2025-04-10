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

	"kafka-notification-system/cmd/producer-service/internal/handler"
	"kafka-notification-system/cmd/producer-service/internal/service"
	"kafka-notification-system/pkg/config"
	"kafka-notification-system/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Producer Service API
// @version 1.0
// @description The Producer Service API for Kafka Notification System
// @host localhost:3000
// @BasePath /
func main() {
	// Загружаем конфигурацию
	appConfig := config.LoadAppConfig()
	kafkaConfig := config.LoadKafkaConfig("producer-service", "")

	// Инициализируем логгер
	logger.InitLogger(appConfig.Environment)
	log := logger.GetLogger()

	// Создаем сервисы
	kafkaService := service.NewKafkaService(kafkaConfig)
	defer func() {
		if err := kafkaService.Close(); err != nil {
			log.Error("Failed to close Kafka service", zap.Error(err))
		}
	}()

	// Создаем обработчики
	producerHandler := handler.NewProducerHandler(kafkaService, log)

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
		v1.POST("/messages", producerHandler.SendMessage)
		v1.GET("/health", producerHandler.Health)
	}

	// Swagger документация
	router.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Info("Starting Producer Service",
			zap.String("port", appConfig.Port),
			zap.Strings("kafka_brokers", kafkaConfig.Brokers))

		fmt.Printf("🚀 Swagger running at http://localhost:%s/api/index.html\n", appConfig.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Producer Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Producer Service stopped")
}
