package location

import (
	"context"
	"encoding/json"

	"github.com/tuanta7/k6-demo/services/location/internal/domain"
	"github.com/tuanta7/k6-demo/services/location/pkg/redis"
)

type Cache struct {
	redis redis.Client
}

func NewCache(redis redis.Client) *Cache {
	return &Cache{
		redis: redis,
	}
}

func (c *Cache) Set(ctx context.Context, location *domain.Location) error {
	return c.redis.Set(ctx, location.TripID, location, 0)
}

func (c *Cache) Get(ctx context.Context, tripID string) (*domain.Location, error) {
	data, err := c.redis.Get(ctx, tripID)
	if err != nil {
		return nil, err
	}

	var l domain.Location
	if err = json.Unmarshal(data, &l); err != nil {
		return nil, err
	}

	location := &domain.Location{}
	err = json.Unmarshal(data, location)
	if err != nil {
		return nil, err
	}

	return location, nil
}
