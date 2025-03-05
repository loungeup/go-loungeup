package restasks

import (
	"fmt"
	"time"

	"github.com/jirenius/go-res/resprot"
	"github.com/nats-io/nats.go"
)

type natsRequester interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
}

// Wait for the task with the given RID to complete.
func Wait[T any](requester natsRequester, taskRID string, options ...waitOption) (T, error) {
	const (
		defaultWaitInterval = time.Second
		defaultWaitTimeout  = time.Minute
	)

	config := &waitConfig{
		interval: defaultWaitInterval,
		timeout:  defaultWaitTimeout,
	}
	for _, option := range options {
		option(config)
	}

	var result T

	ticker := time.NewTicker(config.interval)
	defer ticker.Stop()

	timeout := time.After(config.timeout)

	for {
		select {
		case <-ticker.C:
			message, err := requester.Request("get."+taskRID, nil, config.timeout)
			if err != nil {
				return result, fmt.Errorf("could not get task: %w", err)
			}

			model := &taskRESModel{}
			if _, err := resprot.ParseResponse(message.Data).ParseModel(model); err != nil {
				return result, fmt.Errorf("could not parse task model from response: %w", err)
			}

			if model.isRunning() {
				continue
			}

			if errorMessage := model.Error; errorMessage != "" {
				return result, fmt.Errorf("%s", errorMessage)
			}

			if err := model.decodeResult(&result); err != nil {
				return result, err
			}

			return result, nil
		case <-timeout:
			return result, fmt.Errorf("timeout waiting for task to complete")
		}
	}
}

func WithWaitInterval(interval time.Duration) waitOption {
	return func(config *waitConfig) { config.interval = interval }
}

func WithWaitTimeout(timeout time.Duration) waitOption {
	return func(config *waitConfig) { config.timeout = timeout }
}

type waitOption func(config *waitConfig)

type waitConfig struct {
	interval time.Duration
	timeout  time.Duration
}
