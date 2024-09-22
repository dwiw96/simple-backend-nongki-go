package middleware

import (
	"context"
	"crypto/rsa"
	"os"
	"testing"
	"time"

	cfg "simple-backend-nongki-go/config"
	auth "simple-backend-nongki-go/features/auth"
	pg "simple-backend-nongki-go/utils/driver/postgresql"
	generator "simple-backend-nongki-go/utils/generator"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pool *pgxpool.Pool
	ctx  context.Context
)

func TestMain(m *testing.M) {
	os.Setenv("DB_USERNAME", "dwiw")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "nongki_db")

	envConfig := &cfg.EnvConfig{
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}

	pool = pg.ConnectToPg(envConfig)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	os.Exit(m.Run())
}

func createTokenAndKey(t *testing.T) (string, *rsa.PrivateKey) {
	key, err := LoadKey(ctx, pool)
	require.NoError(t, err)
	require.NotNil(t, key)

	firstname := generator.CreateRandomString(5)
	payload := auth.User{
		Fullname:      firstname + " " + generator.CreateRandomString(7),
		Email:         generator.CreateRandomEmail(firstname),
		Address:       generator.CreateRandomString(20),
		Gender:        generator.CreateRandomGender(),
		MaritalStatus: generator.CreateRandomMaritalStatus(),
	}

	token, err := GetToken(payload, key)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token, key
}

func TestGetToken(t *testing.T) {
	token, _ := createTokenAndKey(t)

	t.Log(token)
}

func TestVerifyToken(t *testing.T) {
	token, key := createTokenAndKey(t)

	t.Run("success", func(t *testing.T) {
		res, err := VerifyToken(token, key)
		require.NoError(t, err)
		require.True(t, res)
	})

	t.Run("failed", func(t *testing.T) {
		res, err := VerifyToken(token+"b", key)
		require.Error(t, err)
		require.False(t, res)
	})
}
