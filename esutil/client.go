package esutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/loungeup/go-loungeup/errors"
	"github.com/loungeup/go-loungeup/log"
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

func NewTypedClient(addresses []string, username, password string) (*elasticsearch.TypedClient, error) {
	result, err := elasticsearch.NewTypedClient(newClientConfig(addresses, username, password))
	if err != nil {
		return nil, fmt.Errorf("could not create Elasticsearch client: %w", err)
	}

	if _, err := result.API.Ping().Do(context.Background()); err != nil {
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
		RetryOnError: func(*http.Request, error) bool {
			// Only use the status code of the HTTP response to retry.
			// Reference: https://github.com/elastic/elastic-transport-go/blob/889f85a00260aae70acbe91eeb18011b60ca7ea8/elastictransport/elastictransport.go#L389-L392
			return false
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

		Logger:    &clientLogger{baseLogger: logger},
		Transport: &clientTransport{baseTransport: http.DefaultTransport},
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

	httpAttrs := func() *httpLogAttrs {
		if (response != nil && response.StatusCode >= http.StatusBadRequest) || err != nil {
			return makeHTTPLogAttrsWithBodies(request, response)
		}

		return makeHTTPLogAttrs(request, response)
	}()
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

	return result
}

func makeHTTPLogAttrsWithBodies(request *http.Request, response *http.Response) *httpLogAttrs {
	result := makeHTTPLogAttrs(request, response)

	if body := request.Body; body != nil {
		result.Request = append(result.Request,
			slog.Any("body", convertReaderToLogAttrValue(request.Header.Get("Content-Type"), body)),
		)
	}

	if body := response.Body; body != nil {
		result.Response = append(result.Response,
			slog.Any("body", convertReaderToLogAttrValue(response.Header.Get("Content-Type"), body)),
		)
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

func convertReaderToLogAttrValue(contentType string, reader io.Reader) any {
	buffer := &bytes.Buffer{}
	_, _ = io.Copy(buffer, reader)

	switch {
	case strings.Contains(contentType, "json"):
		return json.RawMessage(buffer.Bytes())
	default:
		return buffer.String()
	}
}

type clientTransport struct{ baseTransport http.RoundTripper }

var _ (http.RoundTripper) = (*clientTransport)(nil)

func (t *clientTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.baseTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if response.Body == nil || response.Body == http.NoBody {
		return response, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return response, fmt.Errorf("could not read response body: %w", err)
	}

	response.Body = io.NopCloser(bytes.NewReader(body)) // Restore.

	searchResponse := &partialSearchResponse{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(searchResponse); err != nil {
		return response, nil // Ignore.
	}

	if err := parseShardFailures(searchResponse.Shards_.Failures); err != nil {
		return response, err
	}

	return response, nil
}

// Reference: https://github.com/elastic/go-elasticsearch/blob/1dda5df4f11fd5f15264279dd6c773f0f97b9536/typedapi/core/search/response.go#L48
type partialSearchResponse struct {
	Shards_ types.ShardStatistics `json:"_shards"`
}

func parseShardFailures(failures []types.ShardFailure) error {
	if len(failures) == 0 {
		return nil
	}

	failure := failures[0].Reason // We only care about the first failure.

	underlyingError := func() error {
		if reason := failure.Reason; reason != nil {
			return fmt.Errorf("could not execute request because of a shard failure: %s", *reason)
		}

		return fmt.Errorf("could not execute request because of a shard failure")
	}()

	switch failure.Type {
	case "query_shard_exception":
		return &errors.Error{Code: errors.CodeInvalid, UnderlyingError: underlyingError}
	default:
		return underlyingError
	}
}
