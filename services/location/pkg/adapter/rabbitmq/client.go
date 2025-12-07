package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan bool
}

func NewClient(url string) (Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &client{
		conn:    conn,
		channel: channel,
		done:    make(chan bool),
	}, err
}

func (c *client) Close() error {
	close(c.done)

	if err := c.conn.Close(); err != nil {
		return err
	}

	if err := c.channel.Close(); err != nil {
		return err
	}

	return nil
}

func (c *client) DeclareExchange(props ExchangeProperties) error {
	return c.channel.ExchangeDeclare(
		props.Name,
		props.Kind,
		props.Durable,
		props.AutoDelete,
		props.Internal,
		props.NoWait,
		amqp.Table(props.Args),
	)
}

func (c *client) DeclareQueue(props QueueProperties) (Queue, error) {
	q, err := c.channel.QueueDeclare(
		props.Name,
		props.Durable,
		props.AutoDelete,
		props.Exclusive,
		props.NoWait,
		amqp.Table(props.Args),
	)
	if err != nil {
		return Queue{}, err
	}

	return Queue{
		Name:      q.Name,
		Messages:  q.Messages,
		Consumers: q.Consumers,
	}, nil
}

func (c *client) BindQueue(props BindingProperties) error {
	return c.channel.QueueBind(
		props.Queue,
		props.Key,
		props.Exchange,
		props.NoWait,
		amqp.Table(props.Args),
	)
}

func (c *client) Publish(ctx context.Context, exchange, key string, mandatory bool, body []byte) error {
	return c.channel.PublishWithContext(ctx, exchange, key, mandatory,
		false, // deprecated, always false
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

func (c *client) PublishWithConfirm(ctx context.Context, exchange, key string, mandatory bool, body []byte) error {
	if err := c.channel.Confirm(false); err != nil {
		return fmt.Errorf("failed to put channel in confirm mode: %w", err)
	}

	confirms := c.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := c.channel.PublishWithContext(ctx, exchange, key, mandatory,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case confirm, ok := <-confirms:
		if !ok {
			return fmt.Errorf("confirm channel closed")
		}
		if !confirm.Ack {
			return fmt.Errorf("message was nacked by server")
		}
	}

	return nil
}

func (c *client) Consume(
	ctx context.Context,
	queue, consumer string,
	autoAck, exclusive, noWait bool,
	handler ConsumerHandler,
) error {
	deliveries, err := c.channel.ConsumeWithContext(
		ctx,
		queue, consumer,
		autoAck, exclusive, false, noWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.done:
			return nil
		case msg, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}

			if err := handler(ctx, msg.Body); err != nil {
				if !autoAck {
					// Nack with requeue on handler error
					_ = msg.Nack(false, true)
				}
			} else {
				if !autoAck {
					_ = msg.Ack(false)
				}
			}
		}
	}
}
