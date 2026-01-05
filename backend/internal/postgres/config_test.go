package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DSN(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "basic config",
			config: Config{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "password with special characters",
			config: Config{
				Host:     "localhost",
				Port:     5432,
				User:     "admin",
				Password: "p@ss:word/with?special&chars",
				DBName:   "mydb",
				SSLMode:  "require",
			},
			expected: "postgres://admin:p%40ss%3Aword%2Fwith%3Fspecial&chars@localhost:5432/mydb?sslmode=require",
		},
		{
			name: "custom port",
			config: Config{
				Host:     "db.example.com",
				Port:     5433,
				User:     "appuser",
				Password: "secret",
				DBName:   "production",
				SSLMode:  "verify-full",
			},
			expected: "postgres://appuser:secret@db.example.com:5433/production?sslmode=verify-full",
		},
		{
			name: "user with special characters",
			config: Config{
				Host:     "localhost",
				Port:     5432,
				User:     "user@domain.com",
				Password: "pass",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "postgres://user%40domain.com:pass@localhost:5432/testdb?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.DSN()
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestConfig_Defaults(t *testing.T) {
	cfg := Config{}

	assert.Equal(t, "", cfg.Host)
	assert.Equal(t, 0, cfg.Port)
	assert.Equal(t, "", cfg.User)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, "", cfg.DBName)
	assert.Equal(t, "", cfg.SSLMode)
	assert.Equal(t, int32(0), cfg.MaxConns)
	assert.Equal(t, int32(0), cfg.MinConns)
	assert.Equal(t, time.Duration(0), cfg.MaxConnLifetime)
	assert.Equal(t, time.Duration(0), cfg.MaxConnIdleTime)
}
