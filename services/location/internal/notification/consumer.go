package notification

import (
	"context"

	"github.com/tuanta7/k6-demo/services/location/internal"
	"github.com/tuanta7/k6-demo/services/location/pkg/amqp"
)

type Consumer struct {
	rabbitmq    amqp.Consumer
	pushHandler *PushNotificationHandler
}

func NewConsumer(rabbitmq amqp.Consumer, ph *PushNotificationHandler) *Consumer {
	return &Consumer{
		rabbitmq:    rabbitmq,
		pushHandler: ph,
	}
}

func (c *Consumer) ConsumeNotificationQueue(ctx context.Context) error {
	return c.rabbitmq.Consume(ctx, internal.NotificationQueue, "", false, false, c.pushHandler)
}
