// Package jetstreamutil provides utilities for working with NATS JetStream.
package jetstreamutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const defaultStreamName = "default"

func GetOrCreateDefaultStream(js jetstream.JetStream, ctx context.Context) (jetstream.Stream, error) {
	result, err := js.Stream(ctx, defaultStreamName)
	if errors.Is(err, jetstream.ErrStreamNotFound) {
		return createDefaultStream(js, ctx)
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func createDefaultStream(js jetstream.JetStream, ctx context.Context) (jetstream.Stream, error) {
	const (
		streamSubject    = "stream.default.>"
		streamDuplicates = 100 * time.Millisecond
	)

	return js.CreateStream(ctx, jetstream.StreamConfig{
		Name:       defaultStreamName,
		Subjects:   []string{streamSubject},
		Retention:  jetstream.WorkQueuePolicy,
		Duplicates: streamDuplicates,
	})
}

// NewConsumerConfigWithDefaults creates a new consumer config with default values.
func NewConsumerConfigWithDefaults(name, subject string) jetstream.ConsumerConfig {
	//nolint:gomnd,mnd
	return jetstream.ConsumerConfig{
		Name:          name,
		Durable:       name,
		AckWait:       30 * time.Second,
		MaxDeliver:    10,
		FilterSubject: subject,
	}
}

const defaultConsumeInterval = 100 * time.Millisecond

type consumeConfig struct {
	interval time.Duration
}

type consumeOption func(*consumeConfig)

// Consume messages from a JetStream consumer with the given function. The function is responsible for acknowledging or
// rejecting the message.
func Consume(consumer jetstream.Consumer, consumeFunc func(jetstream.Msg), options ...consumeOption) error {
	config := &consumeConfig{
		interval: defaultConsumeInterval,
	}
	for _, option := range options {
		option(config)
	}

	ticker := time.NewTicker(config.interval)
	defer ticker.Stop()

	messages, err := consumer.Messages()
	if err != nil {
		return fmt.Errorf("could not get messages: %w", err)
	}
	defer messages.Stop()

	for range ticker.C {
		message, err := messages.Next()
		if err != nil {
			return fmt.Errorf("could not get next message: %w", err)
		}

		consumeFunc(message)
	}

	return nil
}

// WithConsumeInterval sets the interval at which messages are consumed.
func WithConsumeInterval(interval time.Duration) consumeOption {
	return func(config *consumeConfig) { config.interval = interval }
}
