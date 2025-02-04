package jetstreamutil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/log"
	"github.com/loungeup/go-loungeup/pkg/resutil"
	"github.com/nats-io/nats.go/jetstream"
)

type Throttler struct {
	interval      time.Duration
	logger        *log.Logger
	throttledMsgs sync.Map
}

func NewThrottler(options ...throttlerOption) *Throttler {
	const defaultInterval = time.Second

	result := &Throttler{
		interval: defaultInterval,
		logger:   log.Default(),
	}
	for _, option := range options {
		option(result)
	}

	result.logger.Debug("Throttler created", slog.String("interval", result.interval.String()))

	return result
}

func WithThrottlerInterval(interval time.Duration) throttlerOption {
	return func(throttler *Throttler) { throttler.interval = interval }
}

func WithThrottlerLogger(logger *log.Logger) throttlerOption {
	return func(throttler *Throttler) { throttler.logger = logger }
}

func (throttler *Throttler) Handle(next func(msg jetstream.Msg)) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		key := strings.Join([]string{msg.Subject(), string(msg.Data())}, "-")

		validThrottlerInterval := throttler.parseThrottlerInterval(msg.Data())

		l1 := throttler.logger.With(
			slog.Any("data", json.RawMessage(msg.Data())),
			slog.String("key", key),
			slog.String("subject", msg.Subject()),
			slog.String("traceId", uuid.NewString()),
			slog.String("validThrottlerInterval", validThrottlerInterval.String()),
		)

		if validThrottlerInterval == 0 {
			l1.Debug("Processing message")
			next(msg)
			l1.Debug("Message processed")
			throttler.release(key)

			return
		}

		if throttler.isLocked(key) {
			l1.Debug("Terminating duplicated message")
			msg.Term()

			return
		}

		throttler.lock(key)

		l1.Debug("Throttling message")

		timer := time.NewTimer(validThrottlerInterval)
		ticker := time.NewTicker(validThrottlerInterval / 2) //nolint:all

		go func() {
			defer func() {
				l1.Debug("Message processed")
				timer.Stop()
				ticker.Stop()
				throttler.release(key)
			}()

			for {
				select {
				case <-timer.C:
					l1.Debug("Processing message")
					next(msg)

					return // Terminate the goroutine.
				case <-ticker.C:
					l1.Debug("Message in progress")
					msg.InProgress()
				}
			}
		}()
	}
}

func (throttler *Throttler) parseThrottlerInterval(data []byte) time.Duration {
	request := resutil.NewRequestWithParams(&throttlerParams{})
	_ = json.Unmarshal(data, request)

	if request.Params.throttlerInterval >= 0 && request.Params.throttlerInterval <= time.Minute {
		return request.Params.throttlerInterval
	}

	return throttler.interval
}

func (throttler *Throttler) isLocked(key string) bool {
	_, result := throttler.throttledMsgs.Load(key)
	return result
}

func (throttler *Throttler) lock(key string) { throttler.throttledMsgs.Store(key, time.Now()) }

func (throttler *Throttler) release(key string) { throttler.throttledMsgs.Delete(key) }

type throttlerOption func(throttler *Throttler)

type throttlerParams struct {
	throttlerInterval time.Duration
}

var _ (json.Unmarshaler) = (*throttlerParams)(nil)

func (params *throttlerParams) UnmarshalJSON(data []byte) error {
	type jsonModel struct {
		ThrottlerInterval string `json:"throttlerInterval"`
	}

	model := &jsonModel{}
	if err := json.Unmarshal(data, model); err != nil {
		return fmt.Errorf("could not decode throttler handler params: %w", err)
	}

	if model.ThrottlerInterval == "" {
		*params = throttlerParams{throttlerInterval: -1}

		return nil
	}

	throttlerInterval, err := time.ParseDuration(model.ThrottlerInterval)
	if err != nil {
		return fmt.Errorf("could not parse throttler interval: %w", err)
	}

	*params = throttlerParams{throttlerInterval: throttlerInterval}

	return nil
}
