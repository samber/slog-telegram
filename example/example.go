package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"time"

	"log/slog"

	slogtelegram "github.com/samber/slog-telegram/v2"
)

func main() {
	token := "5977160992:AAGcvh0gwuNQO0tFRy-hKnfvEQux0_CChrw"
	username := "@samuelberthe"

	telegramProxySecret := "secret"
	telegramProxyURL := "https://proxy.example.com/mtproto"

	logger := slog.New(slogtelegram.Option{
		Level:    slog.LevelDebug,
		Token:    token,
		Username: username,
		HTTPClient: &http.Client{
			Transport: NewTransport(telegramProxySecret),
			Timeout:   time.Second * 15,
		},
		APIEndpoint:         telegramProxyURL,
		MessageConfigurator: Configurator,
	}.NewTelegramHandler())
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

// Custom HTTP client

type Transport struct {
	http.RoundTripper
	secret string
}

func NewTransport(secret string) *Transport {
	return &Transport{http.DefaultTransport, secret}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("x-auth-secret", t.secret)
	return t.RoundTripper.RoundTrip(req)
}
