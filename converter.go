package slogtelegram

import (
	"fmt"

	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) string

func DefaultConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) string {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	// handler formatter
	message := fmt.Sprintf("%s\n------------\n\n", record.Message)
	message += attrToTelegramMessage("", attrs)
	return message
}

type MessageConfigurator func(messageConfig tgbotapi.MessageConfig, loggerAttr []slog.Attr) tgbotapi.MessageConfig

func attrToTelegramMessage(base string, attrs []slog.Attr) string {
	message := ""

	for i := range attrs {
		attr := attrs[i]
		k := base + attr.Key
		v := attr.Value
		kind := attr.Value.Kind()

		if kind == slog.KindGroup {
			message += attrToTelegramMessage(k+".", v.Group())
		} else {
			message += fmt.Sprintf("%s: %s\n", k, slogcommon.ValueToString(v))
		}
	}

	return message
}
