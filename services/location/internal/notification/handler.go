package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tuanta7/k6-demo/services/location/internal"
	"github.com/tuanta7/k6-demo/services/location/internal/domain"
	"github.com/tuanta7/k6-demo/services/location/pkg/amqp"
)

type Handler struct {
	consumer amqp.Consumer
}

func NewHandler(consumer amqp.Consumer) *Handler {
	return &Handler{
		consumer: consumer,
	}
}

func (h *Handler) ConsumeNotification(ctx context.Context) error {
	handler := func(ctx context.Context, msg []byte) error {
		var m domain.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			return err
		}

		fmt.Println(m.Channel)
		return nil
	}

	return h.consumer.Consume(ctx, internal.NotificationQueue, "", false, false, handler)
}
