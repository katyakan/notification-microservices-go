package service

import (
	"context"
	"kafka-notification-system/pkg/config"
	"kafka-notification-system/pkg/shared"
	"testing"
	"time"
)

func TestNewKafkaService(t *testing.T) {
	kafkaConfig := &config.KafkaConfig{
		Brokers:            []string{"localhost:9092"},
		ClientID:           "test-producer",
		NotificationsTopic: "test-notifications",
	}

	service := NewKafkaService(kafkaConfig)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if service.config != kafkaConfig {
		t.Error("Expected config to be set")
	}

	if service.writer == nil {
		t.Error("Expected writer to be initialized")
	}

	if service.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestKafkaService_SendMessage(t *testing.T) {
	// This test requires a running Kafka instance
	// Skip if not available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	kafkaConfig := &config.KafkaConfig{
		Brokers:            []string{"localhost:9092"},
		ClientID:           "test-producer",
		NotificationsTopic: "test-notifications",
	}

	service := NewKafkaService(kafkaConfig)
	defer service.Close()

	req := &shared.CreateMessageRequest{
		Type: "notification",
		Payload: map[string]interface{}{
			"chatId": int64(123456),
			"text":   "Test message",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := service.SendMessage(ctx, req)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if response.ID == "" {
		t.Error("Expected non-empty message ID")
	}
}

func TestKafkaService_SendMessage_InvalidPayload(t *testing.T) {
	kafkaConfig := &config.KafkaConfig{
		Brokers:            []string{"localhost:9092"},
		ClientID:           "test-producer",
		NotificationsTopic: "test-notifications",
	}

	service := NewKafkaService(kafkaConfig)
	defer service.Close()

	// Test with invalid payload that can't be marshaled to JSON
	req := &shared.CreateMessageRequest{
		Type:    "notification",
		Payload: make(chan int), // channels can't be marshaled to JSON
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := service.SendMessage(ctx, req)
	if err == nil {
		t.Error("Expected error when sending message with invalid payload")
	}
}
