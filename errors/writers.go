package errors

import (
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type errorLogger interface {
	Error(message string, attributes ...slog.Attr)
	FormattedError(message string, attributes ...slog.Attr)
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

	logAttributes := append(logContext.Attributes(), extractLogAttributes(w)...)

	if err, ok := err.(*res.Error); ok {
		logAttributes = append(logAttributes, slog.Any("underlyingError", err.Data))

		switch err.Code {
		case res.CodeInternalError:
			l.FormattedError(err.Message, logAttributes...)
		default:
			l.Error(err.Message, logAttributes...)
		}

		w.Error(err)

		return
	}

	switch ErrorCode(err) {
	case CodeInternal:
		l.FormattedError(err.Error(), logAttributes...)
	default:
		l.Error(err.Error(), logAttributes...)
	}

	w.Error(&res.Error{Code: getRESErrorCode(err), Message: ErrorMessage(err), Data: logContext})
}

// getRESErrorCode returns the RES error code for a given error.
func getRESErrorCode(err error) string {
	code := ErrorCode(err)

	switch {
	case strings.EqualFold(code, ""):
		return ""
	case strings.EqualFold(code, CodeConflict):
		return res.CodeInvalidParams
	case strings.EqualFold(code, CodeInternal):
		return res.CodeInternalError
	case strings.EqualFold(code, CodeInvalid):
		return res.CodeInvalidParams
	case strings.EqualFold(code, CodeNotFound):
		return res.CodeNotFound
	default:
		return code
	}
}

// extractLogAttributes from the given value.
func extractLogAttributes(value any) []slog.Attr {
	requestAttributes := []any{}

	if request, ok := value.(res.Resource); ok {
		requestAttributes = append(requestAttributes,
			slog.String("name", request.ResourceName()),
			slog.String("query", request.Query()),
		)
	}

	if request, ok := value.(hasType); ok {
		requestAttributes = append(requestAttributes, slog.String("type", request.Type()))
	}

	if request, ok := value.(res.CallRequest); ok {
		requestAttributes = append(requestAttributes,
			slog.Any("params", request.RawParams()),
			slog.Any("token", request.RawToken()),
			slog.Bool("isHttp", request.IsHTTP()),
			slog.String("connectionId", request.CID()),
			slog.String("method", request.Method()),
		)
	}

	if len(requestAttributes) == 0 {
		return nil
	}

	return []slog.Attr{
		slog.Group("request", requestAttributes...),
	}
}

// https://pkg.go.dev/github.com/jirenius/go-res@v0.5.0#Request.Type
type hasType interface {
	Type() string
}
