package cache

import (
	"context"
	"fmt"
	"sync"

	"e-wallet-system-api/config"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
}

var (
	redisOnce     sync.Once
	redisInstance *RedisDatabase
)

func NewRedisDatabase(conf *config.Config) Cache {
	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
			Password: conf.Redis.Password,
			DB:       0,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		redisInstance = &RedisDatabase{client: client}
	})
	return redisInstance
}

func (r *RedisDatabase) GetInstance() *redis.Client {
	return r.client
}
