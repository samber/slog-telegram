package main

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"log/slog"

	slogtelegram "github.com/samber/slog-telegram/v2"
)

// run with
// TOKEN=<your token> CHAT_ID=<your chat id> go run example.go

func main() {
	token := os.Getenv("TOKEN")
	chatId := os.Getenv("CHAT_ID")

	logger := slog.New(slogtelegram.Option{Level: slog.LevelDebug, Token: token, ChatId: chatId, MessageConfigurator: Configurator}.NewTelegramHandler())
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

	// as its async, wait for the message to be sent
	time.Sleep(5 * time.Second)
}

// Configurator Make the message support markdown
func Configurator(config tgbotapi.MessageConfig, attr []slog.Attr) tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	return config
}
