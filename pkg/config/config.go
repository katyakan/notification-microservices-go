package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig содержит общую конфигурацию приложения
type AppConfig struct {
	Port             string `mapstructure:"port"`
	TelegramBotToken string `mapstructure:"telegram_bot_token"`
	Environment      string `mapstructure:"environment"`
}

// LoadAppConfig загружает конфигурацию приложения
func LoadAppConfig() *AppConfig {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Устанавливаем значения по умолчанию
	viper.SetDefault("port", "3000")
	viper.SetDefault("environment", "development")

	// Пытаемся прочитать .env файл (не критично если его нет)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read .env file: %v", err)
	}

	config := &AppConfig{
		Port:             viper.GetString("port"),
		TelegramBotToken: viper.GetString("telegram_bot_token"),
		Environment:      viper.GetString("environment"),
	}

	return config
}
