package handler

import (
	"net/http"

	"kafka-notification-system/pkg/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NotificationHandler обрабатывает HTTP запросы для Notification Service
type NotificationHandler struct {
	logger *zap.Logger
}

// NewNotificationHandler создает новый экземпляр NotificationHandler
func NewNotificationHandler(logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		logger: logger,
	}
}

// Health godoc
// @Summary Health check
// @Description Get the health status of the notification service
// @Tags Health
// @Produce json
// @Success 200 {object} shared.HealthResponse
// @Router /health [get]
func (h *NotificationHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, shared.HealthResponse{Status: "ok"})
}
