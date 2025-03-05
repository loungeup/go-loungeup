package jetstreamutil_test

import (
	"context"
	"time"

	"github.com/loungeup/go-loungeup/jetstreamutil"
	"github.com/nats-io/nats.go/jetstream"
)

func Example() {
	js, err := jetstream.New(nil) // Create a JetStream context as usual.
	if err != nil {
		panic(err)
	}

	// Get (or create) the default stream.
	stream, err := jetstreamutil.GetOrCreateDefaultStream(js, context.Background())
	if err != nil {
		panic(err)
	}

	// Create a pre-configured consumer.
	consumer, err := stream.CreateOrUpdateConsumer(
		context.Background(),
		jetstreamutil.NewConsumerConfigWithDefaults("test", "stream.test"),
	)
	if err != nil {
		panic(err)
	}

	// Consume messages from the consumer.
	if err := jetstreamutil.Consume(consumer, func(message jetstream.Msg) {
		message.Ack() // Process the message.
	}, jetstreamutil.WithConsumeInterval(time.Second)); err != nil {
		panic(err)
	}
}
