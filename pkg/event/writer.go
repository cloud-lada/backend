// Package event provides types for communicating with various event buses.
package event

import (
	"context"
	"encoding/json"
	"time"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
	_ "gocloud.dev/pubsub/azuresb"
	_ "gocloud.dev/pubsub/gcppubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
	_ "gocloud.dev/pubsub/mempubsub"
	_ "gocloud.dev/pubsub/natspubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

type (
	// The Writer type is used to publish JSON-encoded messages onto an event bus.
	Writer struct {
		topic *pubsub.Topic
	}
)

// NewWriter returns a new instance of the Writer type that can be used to publish JSON-encoded messages onto
// an event bus. The event bus used is expressed via the url string.
func NewWriter(ctx context.Context, url string) (*Writer, error) {
	topic, err := pubsub.OpenTopic(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Writer{topic: topic}, nil
}

// Write a JSON-encoded message onto the bus.
func (w *Writer) Write(ctx context.Context, message json.RawMessage) error {
	return w.topic.Send(ctx, &pubsub.Message{
		Body: message,
	})
}

// Close the connection to the event bus.
func (w *Writer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return w.topic.Shutdown(ctx)
}
