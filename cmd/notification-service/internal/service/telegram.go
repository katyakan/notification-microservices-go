package service

import (
	"fmt"
	"kafka-notification-system/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// TelegramService обрабатывает отправку сообщений в Telegram
type TelegramService struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
}

// NewTelegramService создает новый экземпляр TelegramService
func NewTelegramService(botToken string) (*TelegramService, error) {
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	log := logger.GetLogger()
	log.Info("Telegram bot started", zap.String("username", bot.Self.UserName))

	return &TelegramService{
		bot:    bot,
		logger: log,
	}, nil
}

// SendMessage отправляет сообщение в Telegram чат
func (s *TelegramService) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	
	_, err := s.bot.Send(msg)
	if err != nil {
		s.logger.Error("Error sending message to Telegram",
			zap.Error(err),
			zap.Int64("chatId", chatID),
			zap.String("text", text))
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	s.logger.Info("Message sent to Telegram",
		zap.Int64("chatId", chatID),
		zap.String("text", text))

	return nil
}
