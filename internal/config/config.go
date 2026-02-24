package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig
	Store     StoreConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// StoreConfig holds in-memory store settings.
type StoreConfig struct {
	ShardCount int
	GCInterval time.Duration
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	RequestsPerWindow int
	WindowSize        time.Duration
}

// Load reads configuration from environment variables
// and returns a Config with sensible defaults.
func Load() (*Config, error) {

	cfg := &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Store: StoreConfig{
			ShardCount: 64,
			GCInterval: 30 * time.Second,
		},
		RateLimit: RateLimitConfig{
			RequestsPerWindow: 100,
			WindowSize:        60 * time.Second,
		},
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		cfg.Server.Port = port
	}

	if v := os.Getenv("RATE_LIMIT"); v != "" {
		rateLimit, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		cfg.RateLimit.RequestsPerWindow = rateLimit
	}

	if v := os.Getenv("WINDOW_SIZE"); v != "" {
		duration, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
		cfg.RateLimit.WindowSize = duration
	}
	return cfg, nil
}
