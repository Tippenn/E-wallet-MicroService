package messaging

import "context"

// Handler processes one message body. Return nil -> Ack, error -> Nack.
type Handler func(ctx context.Context, body []byte) error

type MessageBroker interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
	Consume(ctx context.Context, queue string, handle Handler) error
	Close() error
}
