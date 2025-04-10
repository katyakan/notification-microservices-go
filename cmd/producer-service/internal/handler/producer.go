package handler

import (
	"context"
	"net/http"

	"kafka-notification-system/pkg/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// KafkaServiceInterface определяет интерфейс для Kafka сервиса
type KafkaServiceInterface interface {
	SendMessage(ctx context.Context, req *shared.CreateMessageRequest) (*shared.CreateMessageResponse, error)
	Close() error
}

// ProducerHandler обрабатывает HTTP запросы для Producer Service
type ProducerHandler struct {
	kafkaService KafkaServiceInterface
	logger       *zap.Logger
}

// NewProducerHandler создает новый экземпляр ProducerHandler
func NewProducerHandler(kafkaService KafkaServiceInterface, logger *zap.Logger) *ProducerHandler {
	return &ProducerHandler{
		kafkaService: kafkaService,
		logger:       logger,
	}
}

// SendMessage godoc
// @Summary Send message to Kafka
// @Description Send a message to Kafka topic
// @Tags Producer
// @Accept json
// @Produce json
// @Param message body shared.CreateMessageRequest true "Message to send"
// @Success 201 {object} shared.CreateMessageResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /messages [post]
func (h *ProducerHandler) SendMessage(c *gin.Context) {
	var req shared.CreateMessageRequest

	// Валидируем входящий JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}

	// Отправляем сообщение в Kafka
	response, err := h.kafkaService.SendMessage(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to send message", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Health godoc
// @Summary Health check
// @Description Get the health status of the producer service
// @Tags Health
// @Produce json
// @Success 200 {object} shared.HealthResponse
// @Router /health [get]
func (h *ProducerHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, shared.HealthResponse{Status: "ok"})
}
