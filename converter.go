package slogtelegram

import (
	"encoding"
	"fmt"
	"strconv"

	"golang.org/x/exp/slog"
)

type Converter func(loggerAttr []slog.Attr, record *slog.Record) string

func DefaultConverter(loggerAttr []slog.Attr, record *slog.Record) string {
	message := fmt.Sprintf("%s\n------------\n\n", record.Message)

	message += attrToTelegramMessage("", loggerAttr)
	record.Attrs(func(attr slog.Attr) bool {
		message += attrToTelegramMessage("", []slog.Attr{attr})
		return true
	})

	return message
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
			message += fmt.Sprintf("%s: %s\n", k, attrToValue(v))
		}
	}

	return message
}

func attrToValue(v slog.Value) string {
	kind := v.Kind()

	switch kind {
	case slog.KindAny:
		return anyValueToString(v)
	case slog.KindLogValuer:
		return anyValueToString(v)
	case slog.KindGroup:
		// not expected to reach this line
		return anyValueToString(v)
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().UTC().String()
	default:
		return anyValueToString(v)
	}
}

func anyValueToString(v slog.Value) string {
	if tm, ok := v.Any().(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return ""
		}

		return string(data)
	}

	return fmt.Sprintf("%+v", v.Any())
}
