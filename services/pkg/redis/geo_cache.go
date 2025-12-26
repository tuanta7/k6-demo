package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type GeoCache interface {
	GeoAdd(ctx context.Context, key, member string, longitude, latitude float64) error
	GeoPos(ctx context.Context, key string, members ...string) ([]*goredis.GeoPos, error)
	GeoSearch(ctx context.Context, longitude, latitude float64, radius float64) ([]*goredis.GeoLocation, error)
}

func (c *Client) GeoAdd(ctx context.Context, key, member string, longitude, latitude float64) error {
	return c.redisClient.GeoAdd(ctx, key, &goredis.GeoLocation{
		Name:      member,
		Longitude: longitude,
		Latitude:  latitude,
	}).Err()
}

func (c *Client) GeoPos(ctx context.Context, key string, members ...string) ([]*goredis.GeoPos, error) {
	result := c.redisClient.GeoPos(ctx, key, members...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	if len(result.Val()) < 1 {
		return nil, goredis.Nil
	}

	return result.Result()
}

func (c *Client) GeoSearch(ctx context.Context, key string, longitude, latitude, radius float64) ([]goredis.GeoLocation, error) {
	result := c.redisClient.GeoSearchLocation(ctx, key, &goredis.GeoSearchLocationQuery{
		GeoSearchQuery: goredis.GeoSearchQuery{
			Radius:    radius,
			Longitude: longitude,
			Latitude:  latitude,
		},
		WithHash: true, // GeoLocation's GeoHash is used here!
	})
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Result()
}
