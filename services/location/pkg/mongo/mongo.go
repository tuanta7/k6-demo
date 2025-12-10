package mongo

import "context"

type DB interface {
	Create(ctx context.Context)
}
