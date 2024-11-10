package config

import (
	"github.com/IBM/sarama"
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

// HTTPConfig interface for HTTPConfig
type HTTPConfig interface {
	Address() string
}

// SwaggerConfig interface for SwaggerConfig
type SwaggerConfig interface {
	Address() string
}

// KafkaConsumerConfig interface for KafkaConsumerConfig
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}
