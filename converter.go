package slogsyslog

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/google/uuid"
	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"
var ErrorKeys = []string{"error", "err"}

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) Message

func DefaultConverter(appName, hostname string) Converter {
	return func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) Message {
		return defaultConverter(appName, hostname, addSource, replaceAttr, loggerAttr, groups, record)
	}
}
func defaultConverter(appName, hostname string, addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) Message {
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	attrs = slogcommon.ReplaceError(attrs, ErrorKeys...)
	if addSource {
		attrs = append(attrs, slogcommon.Source(SourceKey, record))
	}
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	message := Message{
		AppName:   appName,
		Hostname:  hostname,
		Priority:  ConvertSlogToSyslogSeverity(record.Level),
		Timestamp: record.Time.UTC(),
		MessageID: uuid.New().String(),
		Message:   []byte(record.Message),
		ProcessID: strconv.Itoa(os.Getpid()),
	}

	for _, attr := range attrs {
		message.AddStructureData("ID", attr.Key, attr.Value.String())
	}

	return message
}
