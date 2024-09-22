package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvConfig(t *testing.T) {
	ans := &EnvConfig{
		SERVER_PORT:   ":8080",
		DB_USERNAME:   "dwiw",
		DB_PASSWORD:   "secret",
		DB_HOST:       "localhost",
		DB_PORT:       "5432",
		DB_NAME:       "nongki_db",
		ABSOLUTE_PATH: "/home/dwiw22/jobs-test-interview/test-job/toonesia/simple-backend-nongki-go",
	}
	res := GetEnvConfig()
	require.NotNil(t, res)
	assert.Equal(t, ans, res)
}
