package nosql

import (
	"context"
	"time"
)

type NoSQLDB interface {
	// expiration -1 means keep alive
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, obj any) (bool, error)
	Del(ctx context.Context, key string) error
}
