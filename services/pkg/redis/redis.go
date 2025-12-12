package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	cc *goredis.ClusterClient
}

func NewClusterClient(ctx context.Context, opts *goredis.FailoverOptions) (*Client, error) {
	c := goredis.NewFailoverClusterClient(opts)
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

func (c *Client) Exists(ctx context.Context, key ...string) (int64, error) {
	result := c.cc.Exists(ctx, key...)
	if err := result.Err(); err != nil {
		return 0, err
	}

	return result.Val(), nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	result := c.cc.Get(ctx, key)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Bytes()
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.cc.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.cc.Del(ctx, keys...).Err()
}

func (c *Client) MGet(ctx context.Context, keys ...string) ([]any, error) {
	result := c.cc.MGet(ctx, keys...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Result()
}

func (c *Client) MSet(ctx context.Context, pairs ...any) error {
	return c.cc.MSet(ctx, pairs...).Err()
}

func (c *Client) Close() error {
	return c.cc.Close()
}
