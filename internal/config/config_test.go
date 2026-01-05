package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
		"SERVER_HOST": os.Getenv("SERVER_HOST"),
		"SERVER_PORT": os.Getenv("SERVER_PORT"),
	}
	t.Cleanup(func() {
		for k, v := range originalEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	})

	t.Run("loads config with all required env vars", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_SSLMODE", "disable")
		os.Setenv("SERVER_HOST", "127.0.0.1")
		os.Setenv("SERVER_PORT", "9000")

		cfg, err := Load()
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, 5432, cfg.Database.Port)
		assert.Equal(t, "testuser", cfg.Database.User)
		assert.Equal(t, "testpass", cfg.Database.Password)
		assert.Equal(t, "testdb", cfg.Database.DBName)
		assert.Equal(t, "disable", cfg.Database.SSLMode)
		assert.Equal(t, "127.0.0.1", cfg.Server.Host)
		assert.Equal(t, 9000, cfg.Server.Port)
	})

	t.Run("uses default values", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")

		cfg, err := Load()
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, 5432, cfg.Database.Port)
		assert.Equal(t, "require", cfg.Database.SSLMode)
		assert.Equal(t, "0.0.0.0", cfg.Server.Host)
		assert.Equal(t, 8080, cfg.Server.Port)
	})

	t.Run("fails when required DB_HOST is missing", func(t *testing.T) {
		os.Unsetenv("DB_HOST")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")

		cfg, err := Load()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("fails when required DB_USER is missing", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Unsetenv("DB_USER")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")

		cfg, err := Load()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("fails when required DB_PASSWORD is missing", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Unsetenv("DB_PASSWORD")
		os.Setenv("DB_NAME", "testdb")

		cfg, err := Load()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("fails when required DB_NAME is missing", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Unsetenv("DB_NAME")

		cfg, err := Load()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestServerDefaults(t *testing.T) {
	s := Server{}
	assert.Equal(t, "", s.Host)
	assert.Equal(t, 0, s.Port)
}
