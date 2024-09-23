package chache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	cfg "simple-backend-nongki-go/config"
	auth "simple-backend-nongki-go/features/auth"
	middleware "simple-backend-nongki-go/middleware"
	pg "simple-backend-nongki-go/utils/driver/postgresql"
	rd "simple-backend-nongki-go/utils/driver/redis"
	generator "simple-backend-nongki-go/utils/generator"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	chacheTest auth.CacheInterface
	pool       *pgxpool.Pool
	client     *redis.Client
	ctx        context.Context
)

func TestMain(m *testing.M) {
	os.Setenv("DB_USERNAME", "dwiw")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "nongki_db")
	os.Setenv("REDIS_HOST", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "secret")

	envConfig := &cfg.EnvConfig{
		DB_USERNAME:    os.Getenv("DB_USERNAME"),
		DB_PASSWORD:    os.Getenv("DB_PASSWORD"),
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		DB_NAME:        os.Getenv("DB_NAME"),
		REDIS_HOST:     os.Getenv("REDIS_HOST"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
	}

	pool = pg.ConnectToPg(envConfig)

	client = rd.ConnectToRedis(envConfig)
	defer client.Close()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	chacheTest = NewAuthCache(client, ctx)

	os.Exit(m.Run())
}

func TestCachingBlockedToken(t *testing.T) {
	key, err := middleware.LoadKey(ctx, pool)
	require.NoError(t, err)
	require.NotNil(t, key)

	user := auth.User{
		ID:            int64(generator.RandomInt(1, 100)),
		Fullname:      generator.CreateRandomString(5) + " " + generator.CreateRandomString(7),
		Email:         generator.CreateRandomEmail(generator.CreateRandomString(5)),
		Address:       generator.CreateRandomString(20),
		Gender:        generator.CreateRandomGender(),
		MaritalStatus: generator.CreateRandomMaritalStatus(),
	}
	token, err := middleware.CreateToken(user, key)
	require.NoError(t, err)
	require.NotZero(t, len(token))
	t.Log("TOKEN:", token)

	payload, err := middleware.ReadToken(token, key)
	require.NoError(t, err)

	err = chacheTest.CachingBlockedToken(*payload)
	require.NoError(t, err)

	res, err := client.Get(ctx, fmt.Sprint(payload.ID)).Result()
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprint(payload.UserID), res)
}
