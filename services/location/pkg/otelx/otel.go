package otelx

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	goredis "github.com/redis/go-redis/v9"
)

func InstrumentRedis(client *goredis.Client) error {
	if err := redisotel.InstrumentMetrics(client); err != nil {
		return err
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		return err
	}

	return nil
}
