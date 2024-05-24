package errors

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type errorLogger interface {
	Error(message string, attributes ...slog.Attr)
}

type errorWriter interface {
	Error(err error) // Error writes the given error.
}

type logContext struct {
	LogID             string `json:"logId"`
	UnderlyingMessage string `json:"underlyingMessage,omitempty"`
}

func (c *logContext) Attributes() []slog.Attr {
	result := []slog.Attr{slog.String("logId", c.LogID)}
	if c.UnderlyingMessage != "" {
		result = append(result, slog.String("underlyingMessage", c.UnderlyingMessage))
	}

	return result
}

// newLogContext creates a new log context.
func newLogContext() *logContext { return &logContext{LogID: uuid.NewString()} }

// LogAndWriteRESError with the given logger and writer.
func LogAndWriteRESError(l errorLogger, w errorWriter, err error) {
	if err == nil {
		return
	}

	logContext := newLogContext()

	if err, ok := err.(*Error); ok && err.UnderlyingError != nil {
		logContext.UnderlyingMessage = err.UnderlyingError.Error()
	}

	l.Error(err.Error(), logContext.Attributes()...)
	w.Error(&res.Error{Code: getRESErrorCode(err), Message: ErrorMessage(err), Data: logContext})
}

// getRESErrorCode returns the RES error code for a given error.
func getRESErrorCode(err error) string {
	switch ErrorCode(err) {
	case CodeConflict:
		return res.CodeInvalidParams
	case CodeInvalid:
		return res.CodeInvalidParams
	case CodeNotFound:
		return res.CodeNotFound
	default:
		return res.CodeInternalError
	}
}
