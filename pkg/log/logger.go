package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
)

// Log attribute keys used by LoungeUp applications.
const (
	formattedMessageKey = "formattedMessage"
	messageKey          = "message"
	statusKey           = "status"
	timestampKey        = "timestamp"
)

var defaultLogger atomic.Pointer[Logger]

//nolint:gochecknoinits
func init() {
	defaultLogger.Store(newDefaultLogger())
}

// Default returns the default logger.
func Default() *Logger { return defaultLogger.Load() }

// Logger used by LoungeUp applications based on the official log/slog package.
type Logger struct {
	// Adapter for external libraries.
	Adapter *Adapter

	underlyingLogger *slog.Logger
}

// Debug logs a debug message with the given attributes.
func (l *Logger) Debug(message string, attributes ...slog.Attr) {
	l.underlyingLogger.LogAttrs(context.TODO(), slog.LevelDebug, message, attributes...)
}

// Error logs an error message with the given attributes.
func (l *Logger) Error(message string, attributes ...slog.Attr) {
	l.underlyingLogger.LogAttrs(context.TODO(), slog.LevelError, message, attributes...)
}

// FormattedDebug logs a debug message with the given attributes and automatically adds a formatted message attribute.
// The formatted message attribute is used to send logs to Datadog.
func (l *Logger) FormattedDebug(message string, attributes ...slog.Attr) {
	l.Debug(message, append(
		attributes,
		slog.String(formattedMessageKey, formatMessage(message)),
	)...)
}

// FormattedError logs an error message with the given attributes and automatically adds a formatted message attribute.
// The formatted message attribute is used to send logs to Datadog.
func (l *Logger) FormattedError(message string, attributes ...slog.Attr) {
	l.Error(message, append(
		attributes,
		slog.String(formattedMessageKey, formatMessage(message)),
	)...)
}

func formatMessage(message string) string {
	return strings.ToLower(strings.ReplaceAll(message, " ", "-"))
}

func newDefaultLogger() *Logger {
	result := &Logger{
		underlyingLogger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: replaceLogAttribute,
		})),
	}
	result.Adapter = &Adapter{underlyingLogger: result}

	return result
}

// replaceLogAttribute to match LoungeUp format.
func replaceLogAttribute(_ []string, attribute slog.Attr) slog.Attr {
	switch attribute.Key {
	case slog.LevelKey:
		return slog.Attr{Key: statusKey, Value: slog.StringValue(strings.ToLower(attribute.Value.String()))}
	case slog.MessageKey:
		return slog.Attr{Key: messageKey, Value: attribute.Value}
	case slog.TimeKey:
		return slog.Attr{Key: timestampKey, Value: attribute.Value}
	}

	return attribute
}
