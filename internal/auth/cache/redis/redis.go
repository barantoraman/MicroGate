package redis

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"time"

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

func (r *redisStore) Delete(ctx context.Context, token string) error {
	_, err := r.client.Del(ctx, token).Result()
	if err != nil {
		return errors.New("failed to delete session")
	}
	return nil
}

func (r *redisStore) Get(ctx context.Context, sessionToken string) (tokenPkg.Token, error) {
	hash := sha256.Sum256([]byte(sessionToken))
	tHash := hash[:]
	userSession, err := r.client.Get(ctx, string(tHash)).Result()
	if err != nil {
		return tokenPkg.Token{}, errors.New("session not found")
	}
	var session tokenPkg.Token
	if err = json.Unmarshal([]byte(userSession), &session); err != nil {
		return tokenPkg.Token{}, errors.New("failed to unmarshal session")
	}
	return session, nil
}

// Set implements contract.Store.
func (r *redisStore) Set(ctx context.Context, sessionToken *tokenPkg.Token) error {
	session, err := json.Marshal(sessionToken)
	if err != nil {
		return errors.New("failed to marshal sessions")
	}
	if err = r.client.Set(ctx, string(sessionToken.Hash), session, time.Minute*60).Err(); err != nil {
		return errors.New("failed to save session to redis")
	}
	return nil
}
