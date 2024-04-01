package slogtelegram

import (
	"fmt"

	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type ParseMode = string

const (
	ParseModeMarkdown ParseMode = "markdown"
	ParseModeHTML     ParseMode = "html"
)

var SourceKey = "source"

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) string

func DefaultConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) string {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)

	// handler formatter
	message := fmt.Sprintf("%s\n------------\n\n", record.Message)
	message += attrToTelegramMessage("", attrs)
	return message
}

type MessageConfig struct {
	ParseMode ParseMode
}

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
