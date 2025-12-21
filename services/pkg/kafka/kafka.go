package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	client *kgo.Client
}

func NewClient(ctx context.Context, seeds []string, topics []string, group string) (*Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) Publish(ctx context.Context, topic string, key, value []byte) error {
	c.client.Produce(ctx, &kgo.Record{
		Key:   key,
		Value: value,
	}, nil)

	return nil
}
