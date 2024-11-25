package esutil

import (
	"fmt"
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
	logAttrs := []slog.Attr{
		slog.Int64("durationInMilliseconds", duration.Milliseconds()),
		slog.String("duration", duration.String()),
		slog.String("endedAt", startedAt.Add(duration).Format(time.RFC3339)),
		slog.String("startedAt", startedAt.Format(time.RFC3339)),
	}

	if request != nil {
		logAttrs = append(logAttrs, slog.Group("request",
			slog.String("method", request.Method),
			slog.String("url", request.URL.String()),
		))
	}

	if response != nil {
		logAttrs = append(logAttrs, slog.Group("response",
			slog.Int("statusCode", response.StatusCode),
		))
	}

	if err != nil {
		logAttrs = append(logAttrs, slog.Any("error", err))
	}

	switch {
	case err != nil, response.StatusCode >= http.StatusInternalServerError:
		l.baseLogger.Error("Could not execute request", logAttrs...)
	case response != nil && response.StatusCode >= http.StatusBadRequest:
		l.baseLogger.Debug("Could not execute invalid request", logAttrs...)
	default:
		l.baseLogger.Debug("Request executed", logAttrs...)
	}

	return nil
}

func (*clientLogger) RequestBodyEnabled() bool  { return false }
func (*clientLogger) ResponseBodyEnabled() bool { return false }
