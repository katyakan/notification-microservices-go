package shared

import (
	"time"

	"github.com/google/uuid"
)

// generateUUID генерирует новый UUID
func generateUUID() string {
	return uuid.New().String()
}

// GetCurrentTimestamp возвращает текущий timestamp в миллисекундах
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
