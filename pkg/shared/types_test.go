package shared

import (
	"encoding/json"
	"testing"
)

func TestNewKafkaMessage(t *testing.T) {
	messageType := "notification"
	payload := map[string]interface{}{
		"chatId": int64(123456),
		"text":   "Test message",
	}

	message := NewKafkaMessage(messageType, payload)

	if message.Type != messageType {
		t.Errorf("Expected type %s, got %s", messageType, message.Type)
	}

	if message.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if message.Timestamp == 0 {
		t.Error("Expected non-zero timestamp")
	}

	if message.Payload == nil {
		t.Error("Expected non-nil payload")
	}
}

func TestKafkaMessage_ToJSON(t *testing.T) {
	message := NewKafkaMessage("test", map[string]string{"key": "value"})

	jsonData, err := message.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Expected non-empty JSON data")
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}
}

func TestFromJSON(t *testing.T) {
	original := NewKafkaMessage("test", map[string]string{"key": "value"})
	jsonData, _ := original.ToJSON()

	parsed, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.ID != original.ID {
		t.Errorf("Expected ID %s, got %s", original.ID, parsed.ID)
	}

	if parsed.Type != original.Type {
		t.Errorf("Expected type %s, got %s", original.Type, parsed.Type)
	}

	if parsed.Timestamp != original.Timestamp {
		t.Errorf("Expected timestamp %d, got %d", original.Timestamp, parsed.Timestamp)
	}
}

func TestKafkaMessage_IsNotificationMessage(t *testing.T) {
	notificationMessage := NewKafkaMessage("notification", map[string]interface{}{})
	if !notificationMessage.IsNotificationMessage() {
		t.Error("Expected notification message to return true")
	}

	otherMessage := NewKafkaMessage("other", map[string]interface{}{})
	if otherMessage.IsNotificationMessage() {
		t.Error("Expected non-notification message to return false")
	}
}

func TestKafkaMessage_GetNotificationPayload(t *testing.T) {
	payload := map[string]interface{}{
		"chatId": int64(123456),
		"text":   "Test notification",
	}

	message := NewKafkaMessage("notification", payload)

	notification, err := message.GetNotificationPayload()
	if err != nil {
		t.Fatalf("Failed to get notification payload: %v", err)
	}

	if notification.ChatID != 123456 {
		t.Errorf("Expected chatId 123456, got %d", notification.ChatID)
	}

	if notification.Text != "Test notification" {
		t.Errorf("Expected text 'Test notification', got %s", notification.Text)
	}
}

func TestIsValidKafkaMessage(t *testing.T) {
	validMessage := map[string]interface{}{
		"id":        "test-id",
		"type":      "test",
		"payload":   map[string]interface{}{},
		"timestamp": int64(1234567890),
	}

	if !IsValidKafkaMessage(validMessage) {
		t.Error("Expected valid message to return true")
	}

	invalidMessage := map[string]interface{}{
		"id":   "test-id",
		"type": "test",
		// missing payload and timestamp
	}

	if IsValidKafkaMessage(invalidMessage) {
		t.Error("Expected invalid message to return false")
	}
}
