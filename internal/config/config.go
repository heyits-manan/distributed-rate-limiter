package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Store     StoreConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type StoreConfig struct {
	ShardCount int
	GCInterval time.Duration
}

type RateLimitConfig struct {
	RequestsPerWindow int
	WindowSize        time.Duration
}

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
			return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
		cfg.Server.Port = port
	}

	if v := os.Getenv("STORE_SHARD_COUNT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid STORE_SHARD_COUNT: %w", err)
		}
		cfg.Store.ShardCount = n
	}

	if v := os.Getenv("STORE_GC_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid STORE_GC_INTERVAL: %w", err)
		}
		cfg.Store.GCInterval = d
	}

	if v := os.Getenv("RATE_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT: %w", err)
		}
		cfg.RateLimit.RequestsPerWindow = n
	}

	if v := os.Getenv("WINDOW_SIZE"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid WINDOW_SIZE: %w", err)
		}
		cfg.RateLimit.WindowSize = d
	}

	return cfg, nil
}
