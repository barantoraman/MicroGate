package cache

import (
	"context"

	cacheContract "github.com/barantoraman/microgate/internal/auth/cache/contract"
	"github.com/barantoraman/microgate/internal/auth/cache/mock"
	"github.com/barantoraman/microgate/internal/auth/cache/redis"
	"github.com/barantoraman/microgate/pkg/config"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
)

const (
	RedisType    = "redis"
	InMemoryType = "inmemory"
)

func GetStore(ctx context.Context, cfg config.AuthServiceConfigurations, logger loggerContract.Logger) cacheContract.Store {
	switch cfg.StoreType {
	case RedisType:
		return redis.NewRedisStore(ctx, cfg, logger)
	case InMemoryType:
		return mock.NewMockStore()
	default:
		return nil
	}
}
