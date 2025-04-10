package service

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-notification-system/pkg/config"
	"kafka-notification-system/pkg/logger"
	"kafka-notification-system/pkg/shared"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaService обрабатывает получение сообщений из Kafka для отправки уведомлений
type KafkaService struct {
	reader           *kafka.Reader
	deadLetterWriter *kafka.Writer
	telegramService  *TelegramService
	config           *config.KafkaConfig
	logger           *zap.Logger
}

// NewKafkaService создает новый экземпляр KafkaService
func NewKafkaService(kafkaConfig *config.KafkaConfig, telegramService *TelegramService) *KafkaService {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaConfig.Brokers,
		Topic:    kafkaConfig.NotificationsTopic,
		GroupID:  kafkaConfig.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	deadLetterWriter := &kafka.Writer{
		Addr:         kafka.TCP(kafkaConfig.Brokers...),
		Topic:        kafkaConfig.DeadLetterTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	return &KafkaService{
		reader:           reader,
		deadLetterWriter: deadLetterWriter,
		telegramService:  telegramService,
		config:           kafkaConfig,
		logger:           logger.GetLogger(),
	}
}

// StartConsuming начинает потребление сообщений из Kafka
func (s *KafkaService) StartConsuming(ctx context.Context) error {
	s.logger.Info("Starting Notification Service Kafka consumer",
		zap.String("topic", s.config.NotificationsTopic),
		zap.String("groupId", s.config.GroupID))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopping Notification Service Kafka consumer")
			return ctx.Err()
		default:
			message, err := s.reader.ReadMessage(ctx)
			if err != nil {
				s.logger.Error("Failed to read message", zap.Error(err))
				continue
			}

			if err := s.processMessage(ctx, message); err != nil {
				s.logger.Error("Failed to process message", zap.Error(err))
				if err := s.handleDeadLetter(ctx, message); err != nil {
					s.logger.Error("Failed to send message to dead letter topic", zap.Error(err))
				}
			}
		}
	}
}

// processMessage обрабатывает полученное сообщение
func (s *KafkaService) processMessage(ctx context.Context, message kafka.Message) error {
	if len(message.Value) == 0 {
		return fmt.Errorf("empty message value")
	}

	s.logger.Info("Received message", zap.ByteString("value", message.Value))

	// Парсим сообщение
	var rawMessage map[string]interface{}
	if err := json.Unmarshal(message.Value, &rawMessage); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Проверяем валидность структуры
	if !shared.IsValidKafkaMessage(rawMessage) {
		return fmt.Errorf("invalid message format")
	}

	// Конвертируем в типизированное сообщение
	kafkaMessage, err := shared.FromJSON(message.Value)
	if err != nil {
		return fmt.Errorf("failed to parse kafka message: %w", err)
	}

	// Обрабатываем только уведомления
	if kafkaMessage.IsNotificationMessage() {
		return s.processNotification(kafkaMessage)
	}

	s.logger.Warn("Received non-notification message", zap.String("type", kafkaMessage.Type))
	return nil
}

// processNotification обрабатывает уведомление
func (s *KafkaService) processNotification(message *shared.KafkaMessage) error {
	s.logger.Info("Processing notification", zap.String("messageId", message.ID))

	notification, err := message.GetNotificationPayload()
	if err != nil {
		return fmt.Errorf("failed to get notification payload: %w", err)
	}

	// Отправляем сообщение в Telegram
	return s.telegramService.SendMessage(notification.ChatID, notification.Text)
}

// handleDeadLetter отправляет сообщение в dead letter topic
func (s *KafkaService) handleDeadLetter(ctx context.Context, originalMessage kafka.Message) error {
	deadLetterMessage := kafka.Message{
		Key:   originalMessage.Key,
		Value: originalMessage.Value,
		Headers: append(originalMessage.Headers, kafka.Header{
			Key:   "error-timestamp",
			Value: []byte(fmt.Sprintf("%d", shared.GetCurrentTimestamp())),
		}),
	}

	return s.deadLetterWriter.WriteMessages(ctx, deadLetterMessage)
}

// Close закрывает соединения с Kafka
func (s *KafkaService) Close() error {
	if err := s.reader.Close(); err != nil {
		s.logger.Error("Failed to close reader", zap.Error(err))
	}
	if err := s.deadLetterWriter.Close(); err != nil {
		s.logger.Error("Failed to close dead letter writer", zap.Error(err))
	}
	return nil
}
