package redis

import (
	"context"
	"log"

	cacheContract "github.com/barantoraman/microgate/internal/auth/cache/contract"
	"github.com/barantoraman/microgate/pkg/config"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
	"github.com/go-redis/redis/v8"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(ctx context.Context, cfg config.AuthServiceConfigurations) cacheContract.Store {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: cfg.RedisPass,
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to ping redis: %s", err.Error())
	}
	return &redisStore{
		client: client,
	}
}

// Delete implements contract.Store.
func (r *redisStore) Delete(ctx context.Context, token string) error {
	panic("unimplemented")
}

// Get implements contract.Store.
func (r *redisStore) Get(ctx context.Context, sessionToken string) (tokenPkg.Token, error) {
	panic("unimplemented")
}

// Set implements contract.Store.
func (r *redisStore) Set(ctx context.Context, sessionToken *tokenPkg.Token) error {
	panic("unimplemented")
}
