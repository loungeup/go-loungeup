package jetstreamutil

import (
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/log"
	"github.com/nats-io/nats.go/jetstream"
)

type Throttler struct {
	interval         time.Duration
	logger           *log.Logger
	progressInterval time.Duration
	throttledMsgs    sync.Map
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

	if result.progressInterval == 0 {
		result.progressInterval = result.interval / 2 //nolint:gomnd,mnd
	}

	result.logger.Debug("Throttler created",
		slog.String("interval", result.interval.String()),
		slog.String("progressInterval", result.progressInterval.String()),
	)

	return result
}

func WithThrottlerInterval(interval time.Duration) throttlerOption {
	return func(t *Throttler) { t.interval = interval }
}

func WithThrottlerLogger(logger *log.Logger) throttlerOption {
	return func(t *Throttler) { t.logger = logger }
}

func WithThrottlerProgressInterval(progressInterval time.Duration) throttlerOption {
	return func(t *Throttler) { t.progressInterval = progressInterval }
}

func (t *Throttler) Handle(next func(msg jetstream.Msg)) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		key := strings.Join([]string{msg.Subject(), string(msg.Data())}, "-")

		l1 := t.logger.With(
			slog.Any("data", json.RawMessage(msg.Data())),
			slog.String("key", key),
			slog.String("subject", msg.Subject()),
			slog.String("traceId", uuid.NewString()),
		)

		if t.isLocked(key) {
			l1.Debug("Terminating duplicated message")
			msg.Term()

			return
		}

		t.lock(key)

		timer := time.NewTimer(t.interval)
		ticker := time.NewTicker(t.progressInterval)

		go func() {
			defer func() {
				l1.Debug("Message processed")
				timer.Stop()
				ticker.Stop()
				t.release(key)
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

func (t *Throttler) isLocked(key string) bool {
	_, result := t.throttledMsgs.Load(key)
	return result
}

func (t *Throttler) lock(key string) { t.throttledMsgs.Store(key, time.Now()) }

func (t *Throttler) release(key string) { t.throttledMsgs.Delete(key) }

type throttlerOption func(*Throttler)
