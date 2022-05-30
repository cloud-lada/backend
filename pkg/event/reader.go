// Package event provides types used to read events from an event bus.
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
	// The Reader type is used to read individual events from an event stream.
	Reader struct {
		subscription *pubsub.Subscription
	}

	// The ReadFunc type is a function that is invoked per-event when using Reader.Read.
	ReadFunc func(ctx context.Context, message json.RawMessage) error
)

// NewReader returns a new instance of the Reader type, configured to use the event bus described in the URL
// string.
func NewReader(ctx context.Context, url string) (*Reader, error) {
	subscription, err := pubsub.OpenSubscription(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Reader{
		subscription: subscription,
	}, nil
}

// Read events from the bus, invoking the ReadFunc for each. This method blocks until the context is cancelled or
// ReadFunc returns a non-nil error. If the ReadFunc returns a non-nil error, a NACK will be performed on the
// message where possible before returning. Otherwise, an ACK is performed and the next event is requested.
func (r *Reader) Read(ctx context.Context, fn ReadFunc) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := r.subscription.Receive(ctx)
			if err != nil {
				return err
			}

			if err = fn(ctx, message.Body); err != nil {
				nack(message)
				return err
			}

			message.Ack()
		}
	}
}

// Close the connection to the event bus.
func (r *Reader) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return r.subscription.Shutdown(ctx)
}

func nack(message *pubsub.Message) {
	if message.Nackable() {
		message.Nack()
	}
}
