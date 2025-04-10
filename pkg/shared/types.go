package shared

import (
	"encoding/json"
	"time"
)

// BaseKafkaMessage представляет базовую структуру Kafka сообщения
type BaseKafkaMessage struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}

// NotificationMessage представляет сообщение для отправки уведомления в Telegram
type NotificationMessage struct {
	ChatID    int64  `json:"chatId"`
	Text      string `json:"text"`
	MessageID string `json:"messageId,omitempty"`
}

// KafkaMessage представляет типизированное Kafka сообщение
type KafkaMessage struct {
	BaseKafkaMessage
}

// CreateMessageRequest представляет запрос на создание сообщения
type CreateMessageRequest struct {
	Type    string      `json:"type" binding:"required" example:"notification"`
	Payload interface{} `json:"payload" binding:"required" example:"{\"chatId\": 123456, \"text\": \"Hello World\"}"`
}

// CreateMessageResponse представляет ответ на создание сообщения
type CreateMessageResponse struct {
	ID string `json:"id"`
}

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status string `json:"status"`
}

// NewKafkaMessage создает новое Kafka сообщение
func NewKafkaMessage(messageType string, payload interface{}) *KafkaMessage {
	return &KafkaMessage{
		BaseKafkaMessage: BaseKafkaMessage{
			ID:        generateUUID(),
			Type:      messageType,
			Payload:   payload,
			Timestamp: time.Now().UnixMilli(),
		},
	}
}

// ToJSON конвертирует сообщение в JSON
func (m *KafkaMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON создает сообщение из JSON
func FromJSON(data []byte) (*KafkaMessage, error) {
	var msg KafkaMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// IsNotificationMessage проверяет, является ли сообщение уведомлением
func (m *KafkaMessage) IsNotificationMessage() bool {
	return m.Type == "notification"
}

// GetNotificationPayload извлекает payload как NotificationMessage
func (m *KafkaMessage) GetNotificationPayload() (*NotificationMessage, error) {
	payloadBytes, err := json.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	var notification NotificationMessage
	err = json.Unmarshal(payloadBytes, &notification)
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

// IsValidKafkaMessage проверяет валидность структуры Kafka сообщения
func IsValidKafkaMessage(data map[string]interface{}) bool {
	requiredFields := []string{"id", "type", "payload", "timestamp"}
	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			return false
		}
	}
	return true
}
