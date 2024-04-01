package slogtelegram

import (
	"context"
	"fmt"
	"net/http"

	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// Telegram bot token
	Token string
	// ChatId is the id of the chat
	ChatId string

	// optional: customize Telegram message builder
	Converter           Converter
	MessageConfigurator MessageConfigurator

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	// optional: customize HTTP client
	HttpClient *http.Client
}

func (o Option) NewTelegramHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Token == "" {
		panic("missing Telegram token")
	}

	if o.ChatId == "" {
		panic("missing Telegram username")
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	if o.HttpClient == nil {
		o.HttpClient = http.DefaultClient
	}

	err := o.checkInit()
	if err != nil {
		fmt.Println("slog-telegram:", err)
		return nil
	}

	return &TelegramHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*TelegramHandler)(nil)

type TelegramHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string

	client *http.Client
}

func (h *TelegramHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *TelegramHandler) Handle(ctx context.Context, record slog.Record) error {
	msg := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)

	go func() {
		// TODO: handle error here. Probably log it to stderr?
		err := h.option.sendMessage(msg)

		if err != nil {
			fmt.Println("slog-telegram:", err)
		}
	}()

	return nil
}

func (h *TelegramHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TelegramHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *TelegramHandler) WithGroup(name string) slog.Handler {
	return &TelegramHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
