package location

import (
	"context"
	"encoding/json"

	"github.com/tuanta7/k6-demo/services/location/internal/domain"
	"github.com/tuanta7/k6-demo/services/location/pkg/adapter/redis"
)

type Cache struct {
	redis redis.Client
}

func NewCache(redis redis.Client) *Cache {
	return &Cache{
		redis: redis,
	}
}

func (c *Cache) UpdateLocation(ctx context.Context, location *domain.Location) error {
	b, err := c.redis.Get(ctx, location.TripID)
	if err != nil {
		return err
	}

	var l domain.Location
	if err = json.Unmarshal(b, &l); err != nil {
		return err
	}

	if location.Timestamp.Before(l.Timestamp) {
		return nil
	}

	return c.redis.Set(ctx, location.TripID, location, 0)
}
