package jetstreamutil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/log"
	"github.com/loungeup/go-loungeup/resutil"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	minInterval = 0
	maxInterval = time.Minute
)

type Throttler struct {
	interval           time.Duration
	inProgressInterval time.Duration
	logger             *log.Logger
	throttledMsgs      sync.Map
}

func NewThrottler(options ...throttlerOption) *Throttler {
	const (
		defaultInterval           = time.Second
		defaultInProgressInterval = 10 * time.Second
	)

	result := &Throttler{
		interval:           defaultInterval,
		inProgressInterval: defaultInProgressInterval,
		logger:             log.Default(),
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

func WithThrottlerInProgressInterval(interval time.Duration) throttlerOption {
	return func(throttler *Throttler) { throttler.inProgressInterval = interval }
}

func WithThrottlerLogger(logger *log.Logger) throttlerOption {
	return func(throttler *Throttler) { throttler.logger = logger }
}

func (throttler *Throttler) Handle(next func(msg jetstream.Msg)) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		key := strings.Join([]string{msg.Subject(), string(msg.Data())}, "-")

		interval := throttler.extractThrottlerInterval(msg.Data())

		l1 := throttler.logger.With(
			slog.Any("data", json.RawMessage(msg.Data())),
			slog.String("key", key),
			slog.String("subject", msg.Subject()),
			slog.String("traceId", uuid.NewString()),
			slog.String("interval", interval.String()),
		)

		if throttler.isLocked(key) {
			l1.Debug("Terminating duplicated message")
			msg.Term()

			return
		}

		throttler.lock(key)

		l1.Debug("Throttling message")

		timer := time.NewTimer(interval)
		ticker := time.NewTicker(throttler.inProgressInterval)

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

// extractThrottlerInterval from the given data or fallback to the interval of the throttler.
func (throttler *Throttler) extractThrottlerInterval(data []byte) time.Duration {
	request := resutil.NewRequestWithParams(&throttlerParams{})
	_ = json.Unmarshal(data, request)

	if request.Params.throttlerInterval >= minInterval && request.Params.throttlerInterval <= maxInterval {
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
