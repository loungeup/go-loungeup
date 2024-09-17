package restasks

import (
	"fmt"
	"time"

	"github.com/jirenius/go-res/resprot"
	"github.com/nats-io/nats.go"
)

const (
	defaultWaiterInterval = time.Second
	defaultWaiterTimeout  = time.Minute
)

type Requester interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
}

// Wait for the task with the given RID to complete.
func Wait[T any](requester Requester, taskRID string, options ...waitOption) (T, error) {
	config := &waitConfig{
		interval: defaultWaiterInterval,
		timeout:  defaultWaiterTimeout,
	}

	for _, option := range options {
		option(config)
	}

	var result T

	ticker := time.NewTicker(defaultWaiterInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			message, err := requester.Request("get."+taskRID, nil, defaultWaiterTimeout)
			if err != nil {
				return result, fmt.Errorf("could not get task: %w", err)
			}

			task := &taskModel{}
			if _, err := resprot.ParseResponse(message.Data).ParseModel(task); err != nil {
				return result, fmt.Errorf("could not parse task model from response: %w", err)
			}

			if task.isRunning() {
				continue
			}

			if err := task.err(); err != nil {
				return result, err
			}

			if err := task.decodeResult(&result); err != nil {
				return result, fmt.Errorf("could not decode task result: %w", err)
			}

			return result, nil
		}
	}
}

func WithWaitInterval(interval time.Duration) waitOption {
	return func(config *waitConfig) { config.interval = interval }
}

func WithWaitTimeout(timeout time.Duration) waitOption {
	return func(config *waitConfig) { config.timeout = timeout }
}

type waitOption func(*waitConfig)

type waitConfig struct {
	interval time.Duration
	timeout  time.Duration
}
