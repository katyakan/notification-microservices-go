package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// KafkaConfig содержит конфигурацию для Kafka
type KafkaConfig struct {
	Brokers              []string      `mapstructure:"brokers"`
	ClientID             string        `mapstructure:"client_id"`
	GroupID              string        `mapstructure:"group_id"`
	RetryInitialTime     time.Duration `mapstructure:"retry_initial_time"`
	RetryMaxAttempts     int           `mapstructure:"retry_max_attempts"`
	NotificationsTopic   string        `mapstructure:"notifications_topic"`
	DeadLetterTopic      string        `mapstructure:"dead_letter_topic"`
}

// LoadKafkaConfig загружает конфигурацию Kafka
func LoadKafkaConfig(clientID, groupID string) *KafkaConfig {
	// Устанавливаем значения по умолчанию
	viper.SetDefault("kafka_brokers", "localhost:9092")
	viper.SetDefault("retry_initial_time", 100*time.Millisecond)
	viper.SetDefault("retry_max_attempts", 8)
	viper.SetDefault("notifications_topic", "notifications")
	viper.SetDefault("dead_letter_topic", "dead-letter")

	// Читаем переменные окружения
	viper.AutomaticEnv()

	brokersStr := viper.GetString("kafka_brokers")
	brokers := strings.Split(brokersStr, ",")

	return &KafkaConfig{
		Brokers:              brokers,
		ClientID:             clientID,
		GroupID:              groupID,
		RetryInitialTime:     viper.GetDuration("retry_initial_time"),
		RetryMaxAttempts:     viper.GetInt("retry_max_attempts"),
		NotificationsTopic:   viper.GetString("notifications_topic"),
		DeadLetterTopic:      viper.GetString("dead_letter_topic"),
	}
}
