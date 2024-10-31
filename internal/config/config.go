package config

import (
	"time"

	"github.com/joho/godotenv"
)

// Load loads environment
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig interface for GRPCConfig
type GRPCConfig interface {
	Address() string
}

// PGConfig interface for PGConfig
type PGConfig interface {
	DSN() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	ElementTTL() time.Duration
}
