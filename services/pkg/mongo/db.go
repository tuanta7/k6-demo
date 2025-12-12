package mongo

import "context"

type DB interface {
	Get(ctx context.Context, collection string, filter any, result any) error
}
