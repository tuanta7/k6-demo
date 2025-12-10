package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	MGet(ctx context.Context, keys ...string) ([]any, error)
	MSet(ctx context.Context, values ...any) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}

type Client struct {
	c *goredis.Client
}

func NewClient(ctx context.Context, opts *goredis.Options) (*Client, error) {
	c := goredis.NewClient(opts)
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

func (c *Client) Exists(ctx context.Context, key ...string) (int64, error) {
	result := c.c.Exists(ctx, key...)
	if err := result.Err(); err != nil {
		return 0, err
	}

	return result.Val(), nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	result := c.c.Get(ctx, key)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Bytes()
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.c.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.c.Del(ctx, keys...).Err()
}

func (c *Client) MGet(ctx context.Context, keys ...string) ([]any, error) {
	result := c.c.MGet(ctx, keys...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Result()
}

func (c *Client) MSet(ctx context.Context, pairs ...any) error {
	return c.c.MSet(ctx, pairs...).Err()
}

func (c *Client) Close() error {
	return c.c.Close()
}
