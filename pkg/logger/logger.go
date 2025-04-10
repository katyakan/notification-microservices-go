package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger инициализирует логгер
func InitLogger(environment string) {
	var err error
	
	if environment == "production" {
		Logger, err = zap.NewProduction()
	} else {
		Logger, err = zap.NewDevelopment()
	}
	
	if err != nil {
		panic(err)
	}
}

// GetLogger возвращает экземпляр логгера
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger("development")
	}
	return Logger
}
