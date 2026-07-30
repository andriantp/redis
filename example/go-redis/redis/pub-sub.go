package redis

import (
	"context"
	"fmt"
)

func (r *repository) PUBLISH(ctx context.Context, channel, message string) error {
	return r.client.Publish(
		ctx,
		channel,
		message,
	).Err()
}

func (r *repository) SUBSCRIBE(ctx context.Context, channel string) error {
	pubsub := r.client.Subscribe(ctx, channel)
	ch := pubsub.Channel()
	fmt.Printf("Subscribed to %q\n", channel)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg := <-ch:
			fmt.Printf(
				"Received message from %q: %s\n",
				msg.Channel,
				msg.Payload,
			)
		}
	}
}
