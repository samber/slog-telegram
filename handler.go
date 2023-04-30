package slogtelegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"
)

// curl -X POST \
//      -H 'Content-Type: application/json' \
//      -d '{"chat_id": "<your-chat-id>", "text": "This is a test from curl", "disable_notification": true}' \
//      https://api.telegram.org/bot<your-bot-token>/sendMessage

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// Telegram bot token
	Token string
	// Username of the channel in the form of `@username`
	Username string

	// optional: customize Telegram message builder
	Converter Converter
}

func (o Option) NewTelegramHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Token == "" {
		panic("missing Telegram token")
	}

	if o.Username == "" {
		panic("missing Telegram username")
	}

	client, err := tgbotapi.NewBotAPI(o.Token)
	if err != nil {
		fmt.Println("slog-telegram:", err)
		return nil
	}

	return &TelegramHandler{
		option: o,
		client: client,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

type TelegramHandler struct {
	option Option
	client *tgbotapi.BotAPI
	attrs  []slog.Attr
	groups []string
}

func (h *TelegramHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *TelegramHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	message := converter(h.attrs, &record)
	msg := tgbotapi.NewMessageToChannel(h.option.Username, message)

	_, err := h.client.Send(msg)
	if err != nil {
		fmt.Println("slog-telegram:", err.Error())
		return err
	}

	return nil
}

func (h *TelegramHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TelegramHandler{
		option: h.option,
		client: h.client,
		attrs:  appendAttrsToGroup(h.groups, h.attrs, attrs),
		groups: h.groups,
	}
}

func (h *TelegramHandler) WithGroup(name string) slog.Handler {
	return &TelegramHandler{
		option: h.option,
		client: h.client,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
