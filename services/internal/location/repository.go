package location

import (
	"context"
	"encoding/json"

	"github.com/tuanta7/k6-demo/services/internal/domain"
	"github.com/tuanta7/k6-demo/services/pkg/redis"
)

type Repository struct {
	redis redis.Cache
}

func NewRepository(cache redis.Cache) *Repository {
	return &Repository{redis: cache}
}
func (c *Repository) Set(ctx context.Context, location *domain.Location) error {
	return c.redis.Set(ctx, location.TripID, location, 0)
}

func (c *Repository) Get(ctx context.Context, tripID string) (*domain.Location, error) {
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
