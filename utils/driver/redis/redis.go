package redisPkg

import (
	"github.com/redis/go-redis/v9"

	cfg "simple-backend-nongki-go/config"
)

func ConnectToRedis(envCfg *cfg.EnvConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     envCfg.REDIS_HOST,
		Password: envCfg.REDIS_PASSWORD, // no password set
		DB:       0,                     // use default DB
	})

	return client
}
