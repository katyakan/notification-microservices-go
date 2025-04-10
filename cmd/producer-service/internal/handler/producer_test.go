package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"kafka-notification-system/pkg/shared"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MockKafkaService для тестирования
type MockKafkaService struct {
	shouldError bool
}

func (m *MockKafkaService) SendMessage(ctx context.Context, req *shared.CreateMessageRequest) (*shared.CreateMessageResponse, error) {
	if m.shouldError {
		return nil, &MockError{message: "mock error"}
	}
	return &shared.CreateMessageResponse{ID: "test-id"}, nil
}

func (m *MockKafkaService) Close() error {
	return nil
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestProducerHandler_SendMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockKafkaService{shouldError: false}
	logger := zap.NewNop()
	handler := NewProducerHandler(mockService, logger)

	router := gin.New()
	router.POST("/messages", handler.SendMessage)

	requestBody := shared.CreateMessageRequest{
		Type: "notification",
		Payload: map[string]interface{}{
			"chatId": int64(123456),
			"text":   "Test message",
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response shared.CreateMessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", response.ID)
	}
}

func TestProducerHandler_SendMessage_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockKafkaService{shouldError: false}
	logger := zap.NewNop()
	handler := NewProducerHandler(mockService, logger)

	router := gin.New()
	router.POST("/messages", handler.SendMessage)

	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestProducerHandler_SendMessage_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockKafkaService{shouldError: true}
	logger := zap.NewNop()
	handler := NewProducerHandler(mockService, logger)

	router := gin.New()
	router.POST("/messages", handler.SendMessage)

	requestBody := shared.CreateMessageRequest{
		Type: "notification",
		Payload: map[string]interface{}{
			"chatId": int64(123456),
			"text":   "Test message",
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestProducerHandler_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockKafkaService{}
	logger := zap.NewNop()
	handler := NewProducerHandler(mockService, logger)

	router := gin.New()
	router.GET("/health", handler.Health)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response shared.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got %s", response.Status)
	}
}
