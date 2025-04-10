package service

import (
	"context"
	"fmt"
	"kafka-notification-system/pkg/config"
	"kafka-notification-system/pkg/logger"
	"kafka-notification-system/pkg/shared"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaService обрабатывает отправку сообщений в Kafka
type KafkaService struct {
	writer *kafka.Writer
	config *config.KafkaConfig
	logger *zap.Logger
}

// NewKafkaService создает новый экземпляр KafkaService
func NewKafkaService(kafkaConfig *config.KafkaConfig) *KafkaService {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaConfig.Brokers...),
		Topic:        kafkaConfig.NotificationsTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	return &KafkaService{
		writer: writer,
		config: kafkaConfig,
		logger: logger.GetLogger(),
	}
}

// SendMessage отправляет сообщение в Kafka
func (s *KafkaService) SendMessage(ctx context.Context, req *shared.CreateMessageRequest) (*shared.CreateMessageResponse, error) {
	// Создаем Kafka сообщение
	message := shared.NewKafkaMessage(req.Type, req.Payload)

	// Конвертируем в JSON
	messageBytes, err := message.ToJSON()
	if err != nil {
		s.logger.Error("Failed to marshal message", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Создаем Kafka сообщение для отправки
	kafkaMessage := kafka.Message{
		Key:   []byte(message.ID),
		Value: messageBytes,
		Headers: []kafka.Header{
			{
				Key:   "message-type",
				Value: []byte(req.Type),
			},
		},
	}

	// Отправляем сообщение
	err = s.writer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		s.logger.Error("Failed to send message to Kafka", 
			zap.Error(err),
			zap.String("messageId", message.ID),
			zap.String("messageType", req.Type))
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	s.logger.Info("Message sent successfully",
		zap.String("messageId", message.ID),
		zap.String("messageType", req.Type))

	return &shared.CreateMessageResponse{ID: message.ID}, nil
}

// Close закрывает соединение с Kafka
func (s *KafkaService) Close() error {
	return s.writer.Close()
}
