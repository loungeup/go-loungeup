package esutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/loungeup/go-loungeup/pkg/log"
)

func NewClient(addresses []string, username, password string) (*elasticsearch.Client, error) {
	result, err := elasticsearch.NewClient(newClientConfig(addresses, username, password))
	if err != nil {
		return nil, fmt.Errorf("could not create Elasticsearch client: %w", err)
	}

	if _, err := result.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping Elasticsearch server: %w", err)
	}

	return result, nil
}

func newClientConfig(addresses []string, username, password string) elasticsearch.Config {
	logger := log.Default().With(slog.String("component", "elasticsearch"))

	//nolint:mnd
	return elasticsearch.Config{
		// Authentication.
		Addresses: addresses,
		Username:  username,
		Password:  password,

		// Retry.
		DisableRetry: false,
		MaxRetries:   5,
		RetryBackoff: func(attempt int) time.Duration {
			result := func() time.Duration {
				switch attempt {
				case 0:
					return 100 * time.Millisecond
				case 1:
					return 200 * time.Millisecond
				case 2:
					return 500 * time.Millisecond
				case 3:
					return time.Second
				case 4:
					return 5 * time.Second
				default:
					return 10 * time.Second
				}
			}()

			logger.Debug("Retrying to execute request",
				slog.Int("attempt", attempt),
				slog.Int64("retryDurationInMilliseconds", result.Milliseconds()),
				slog.String("retryAt", time.Now().Add(result).Format(time.RFC3339)),
				slog.String("retryDuration", result.String()),
			)

			return result
		},
		RetryOnStatus: []int{
			// Client errors. Some might be temporary. Note that 404 is excluded to prevent infinite loops.
			http.StatusBadRequest,
			http.StatusForbidden,
			http.StatusRequestTimeout,
			http.StatusConflict,
			http.StatusTooManyRequests,

			// Server errors.
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},

		Logger: &clientLogger{baseLogger: logger},
	}
}

type clientLogger struct{ baseLogger *log.Logger }

var _ (elastictransport.Logger) = (*clientLogger)(nil)

func (l *clientLogger) LogRoundTrip(
	request *http.Request,
	response *http.Response,
	err error,
	startedAt time.Time,
	duration time.Duration,
) error {
	attrs := []slog.Attr{
		slog.Int64("durationInMilliseconds", duration.Milliseconds()),
		slog.String("duration", duration.String()),
		slog.String("endedAt", startedAt.Add(duration).Format(time.RFC3339)),
		slog.String("startedAt", startedAt.Format(time.RFC3339)),
	}

	httpAttrs := makeHTTPLogAttrs(request, response)
	attrs = append(attrs,
		slog.Group("request", convertLogAttrsToAny(httpAttrs.Request)...),
		slog.Group("response", convertLogAttrsToAny(httpAttrs.Response)...),
	)

	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}

	switch {
	case err != nil, response.StatusCode >= http.StatusInternalServerError:
		l.baseLogger.Error("Could not execute request", attrs...)
	case response != nil && response.StatusCode >= http.StatusBadRequest:
		l.baseLogger.Debug("Could not execute invalid request", attrs...)
	default:
		l.baseLogger.Debug("Request executed", attrs...)
	}

	return nil
}

func (*clientLogger) RequestBodyEnabled() bool  { return true }
func (*clientLogger) ResponseBodyEnabled() bool { return true }

type httpLogAttrs struct {
	Request  []slog.Attr
	Response []slog.Attr
}

func makeHTTPLogAttrs(request *http.Request, response *http.Response) *httpLogAttrs {
	result := &httpLogAttrs{
		Request: []slog.Attr{
			slog.String("method", request.Method),
			slog.String("url", request.URL.String()),
		},
		Response: []slog.Attr{},
	}

	if response == nil {
		return result
	}

	result.Response = append(result.Response, slog.Int("statusCode", response.StatusCode))

	if response.StatusCode >= http.StatusBadRequest {
		if body := request.Body; body != nil {
			result.Request = append(result.Request, slog.Any("body", convertReaderToJSON(body)))
		}

		if body := response.Body; body != nil {
			result.Response = append(result.Response, slog.Any("body", convertReaderToJSON(body)))
		}
	}

	return result
}

func convertLogAttrsToAny(attrs []slog.Attr) []any {
	result := []any{}
	for _, attr := range attrs {
		result = append(result, attr)
	}

	return result
}

func convertReaderToJSON(reader io.Reader) json.RawMessage {
	buffer := &bytes.Buffer{}
	_, _ = io.Copy(buffer, reader)

	return json.RawMessage(buffer.Bytes())
}
