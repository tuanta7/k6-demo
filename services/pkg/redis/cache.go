package redis

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	MGet(ctx context.Context, keys ...string) ([]any, error)
	MSet(ctx context.Context, values ...any) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}
