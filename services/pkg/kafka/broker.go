package kafka

import "context"

type Broker interface{}

type Publisher interface {
	Publish(ctx context.Context, topic string, key, value []byte) error
}

type Consumer interface {
	Consume(ctx context.Context, topic string, handler func([]byte)) error
}
