package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"sync/atomic"

	"github.com/DataDog/gostackparse"
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
	defaultLogger.Store(NewLogger())
}

// Default returns the default logger.
func Default() *Logger { return defaultLogger.Load() }

// Logger used by LoungeUp applications based on the official log/slog package.
type Logger struct{ underlyingLogger *slog.Logger }

// LoggerOption is a type of function that configures a Logger.
type LoggerOption func(*Logger)

// NewLogger creates a new Logger with the given options.
func NewLogger(options ...LoggerOption) *Logger {
	result := &Logger{
		underlyingLogger: newUnderlyingLoggerWithWriter(os.Stdout),
	}
	for _, option := range options {
		option(result)
	}

	return result
}

// WithLoggerWriter sets the writer of the logger.
func WithLoggerWriter(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.underlyingLogger = newUnderlyingLoggerWithWriter(w)
	}
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

// With works like the With method of the official log/slog package.
func (l *Logger) With(attributes ...slog.Attr) *Logger {
	return &Logger{
		underlyingLogger: slog.New(l.underlyingLogger.Handler().WithAttrs(attributes)),
	}
}

// WithGroup works like the WithGroup method of the official log/slog package.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		underlyingLogger: slog.New(l.underlyingLogger.Handler().WithGroup(name)),
	}
}

func HandlePanic() {
	errorValue := recover()
	if errorValue == nil {
		return
	}

	Default().FormattedError("Handling panic",
		slog.Any("errorValue", errorValue),
		slog.Any("goroutinesStack", encodeGoroutinesStack()),
	)

	os.Exit(1)
}

func encodeGoroutinesStack() json.RawMessage {
	goroutinesStack, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))
	result, _ := json.Marshal(goroutinesStack)

	return result
}

func formatMessage(message string) string {
	return strings.ToLower(strings.ReplaceAll(message, " ", "-"))
}

func newUnderlyingLoggerWithWriter(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceLogAttribute,
	}))
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
