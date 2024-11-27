package env

import (
	"net"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	jaegerHostEnvName = "JAEGER_AGENT_HOST"
	jaegerPortEnvName = "JAEGER_AGENT_PORT"

	jaegerSamplerTypeEnvName  = "JAEGER_SAMPLER_TYPE"
	jaegerSamplerParamEnvName = "JAEGER_SAMPLER_PARAM"
	jaegerServiceEnvName      = "JAEGER_SERVICE_NAME"
)

// JaegerConfig interface for jaeger config
type JaegerConfig interface {
	LocalAgentAddress() string
	SamplerType() string
	SamplerParam() float64
	ServiceName() string
}

type jaegerConfig struct {
	host         string
	port         string
	samplerType  string
	serviceName  string
	samplerParam float64
}

// NewJaegerConfig creates new jaeger config
func NewJaegerConfig() (JaegerConfig, error) {
	host := os.Getenv(jaegerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("jaeger agent host not found")
	}

	port := os.Getenv(jaegerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("jaeger agent port not found")
	}

	samplerType := os.Getenv(jaegerSamplerTypeEnvName)
	if len(samplerType) == 0 {
		return nil, errors.New("jaeger sampler type not found")
	}

	samplerParam, err := strconv.ParseFloat(os.Getenv(jaegerSamplerParamEnvName), 64)
	if err != nil {
		return nil, errors.New("jaeger sampler param not found")
	}

	serviceName := os.Getenv(jaegerServiceEnvName)
	if len(serviceName) == 0 {
		return nil, errors.New("jaeger service name not found")
	}

	return &jaegerConfig{
		host:         host,
		port:         port,
		samplerType:  samplerType,
		samplerParam: samplerParam,
		serviceName:  serviceName,
	}, nil
}

// LocalAgentAddress returns address of host
func (cfg *jaegerConfig) LocalAgentAddress() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// SamplerType returns sampler type
func (cfg *jaegerConfig) SamplerType() string {
	return cfg.samplerType
}

// SamplerParam returns sampler param
func (cfg *jaegerConfig) SamplerParam() float64 {
	return cfg.samplerParam
}

// ServiceName returns service name
func (cfg *jaegerConfig) ServiceName() string {
	return cfg.serviceName
}
