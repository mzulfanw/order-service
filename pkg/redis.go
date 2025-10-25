package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/mzulfanw/order-service/configs"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
}

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(config configs.RedisConfig) Cache {
	dsn := fmt.Sprintf("%s:%s", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{Addr: dsn})
	return &redisCache{client: rdb, ctx: context.Background()}
}

func (r *redisCache) Get(key string) ([]byte, bool) {
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (r *redisCache) Set(key string, value []byte) {
	r.client.Set(r.ctx, key, value, 10*time.Minute)
}
