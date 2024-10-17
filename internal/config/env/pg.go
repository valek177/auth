package env

import (
	"errors"
	"os"

	"github.com/valek177/auth/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	dsn string
}

// NewPGConfig creates new postgres config struct
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN returns DSN from config
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
