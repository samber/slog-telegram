package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"log/slog"

	slogtelegram "github.com/samber/slog-telegram/v2"
)

func main() {
	token := "5977160992:AAGcvh0gwuNQO0tFRy-hKnfvEQux0_CChrw"
	username := "@samuelberthe"

	logger := slog.New(slogtelegram.Option{Level: slog.LevelDebug, Token: token, ChatId: username, MessageConfigurator: Configurator}.NewTelegramHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}

// Configurator Make the message support markdown
func Configurator(config tgbotapi.MessageConfig, attr []slog.Attr) tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	return config
}
