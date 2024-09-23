package redisPkg

import (
	"context"
	"os"
	"testing"
	"time"

	config "simple-backend-nongki-go/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectToRedis(t *testing.T) {
	os.Setenv("REDIS_HOST", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "secret")

	envConfig := &config.EnvConfig{
		REDIS_HOST:     os.Getenv("REDIS_HOST"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
	}

	client := ConnectToRedis(envConfig)
	defer client.Close()
	require.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	require.NoError(t, err)

	err = client.Set(ctx, "test", "test redis", 0).Err()
	require.NoError(t, err)

	val, err := client.Get(ctx, "test").Result()
	require.NoError(t, err)
	assert.Equal(t, "test redis", val)
}
