package main

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	slogtelegram "github.com/samber/slog-telegram/v2"
)

// run with
// TOKEN=<your token> CHAT_ID=<your chat id> go run example.go

func main() {
	token := os.Getenv("TOKEN")
	chatId := os.Getenv("CHAT_ID")

	logger := slog.New(slogtelegram.Option{Level: slog.LevelDebug, Token: token, ChatId: chatId, ParseMode: slogtelegram.ParseModeHTML}.NewTelegramHandler())
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
		Error("Hello <b><i>slog</i></b>")

	// as its async, wait for the message to be sent
	time.Sleep(5 * time.Second)
}
