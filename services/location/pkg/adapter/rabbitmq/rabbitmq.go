package rabbitmq

import (
	"context"
)

type Client interface {
	Close() error
	Consumer
	Publisher
}

type Consumer interface {
	Consume(ctx context.Context, queue, consumer string, autoAck, exclusive, noWait bool, handler ConsumerHandler) error
}

type Publisher interface {
	DeclareExchange(props ExchangeProps) error
	DeclareQueue(props QueueProps) (Queue, error)
	BindQueue(props BindingProps) error
	Publish(ctx context.Context, exchange, key string, mandatory bool, body []byte) error
	PublishWithConfirm(ctx context.Context, exchange, key string, mandatory bool, body []byte) error
}

type Arguments map[string]any

type ConsumerHandler func(ctx context.Context, body []byte) error

type ExchangeProps struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       Arguments
}

type QueueProps struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       Arguments
}

type BindingProps struct {
	Queue    string
	Exchange string
	Key      string
	NoWait   bool
	Args     Arguments
}

type Queue struct {
	Name      string
	Messages  int
	Consumers int
}
