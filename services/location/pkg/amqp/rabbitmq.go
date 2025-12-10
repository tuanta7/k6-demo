package amqp

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	QueueDeclare(queue string, durable, autoDelete, exclusive bool, args amqp.Table) (amqp.Queue, error)
	ExchangeDeclare(exchange, kind string, durable, autoDelete bool, args amqp.Table) error
	QueueBind(queue, exchange, key string, args amqp.Table) error
	ExchangeBind(destination, source, key string, args amqp.Table) error
	Close() error
}

type client struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
}

func NewClient(url string) (Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &client{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *client) Close() error {
	if c.channel != nil {
		_ = c.channel.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *client) QueueDeclare(queue string, durable, autoDelete, exclusive bool, args amqp.Table) (amqp.Queue, error) {
	return c.channel.QueueDeclare(queue, durable, autoDelete, exclusive, false, args)
}

func (c *client) ExchangeDeclare(exchange, kind string, durable, autoDelete bool, args amqp.Table) error {
	return c.channel.ExchangeDeclare(exchange, kind, durable, autoDelete, false, false, args)
}

func (c *client) QueueBind(queue, exchange, key string, args amqp.Table) error {
	return c.channel.QueueBind(queue, exchange, key, false, args)
}

func (c *client) ExchangeBind(destination, source, key string, args amqp.Table) error {
	return c.channel.ExchangeBind(destination, source, key, false, args)
}

type Publisher interface {
	Publish(ctx context.Context, exchange, key string, mandatory bool, body []byte) error
	Close() error
}

func (c *client) Publish(ctx context.Context, exchange, key string, mandatory bool, body []byte) error {
	return c.channel.PublishWithContext(ctx, exchange, key, mandatory, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

type ConsumerHandler func(ctx context.Context, body []byte) error

type Consumer interface {
	Consume(ctx context.Context, queue, consumer string, autoAck, exclusive bool, handler ConsumerHandler) error
	Close() error
}
