package handler

import (
	"net/http"

	"kafka-notification-system/pkg/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


type ConsumerHandler struct {
	logger *zap.Logger
}


func NewConsumerHandler(logger *zap.Logger) *ConsumerHandler {
	return &ConsumerHandler{
		logger: logger,
	}
}

// Health godoc
// @Summary Health check
// @Description Get the health status of the consumer service
// @Tags Health
// @Produce json
// @Success 200 {object} shared.HealthResponse
// @Router /health [get]
func (h *ConsumerHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, shared.HealthResponse{Status: "ok"})
}
