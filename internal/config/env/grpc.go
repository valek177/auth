package env

import (
	"net"
	"os"

	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/config"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName    = "GRPC_HOST"
	grpcPortEnvName    = "GRPC_PORT"
	serviceTlsCertFile = "GRPC_TLS_CERT_FILE"
	serviceTlsKeyFile  = "GRPC_TLS_KEY_FILE"
)

type grpcConfig struct {
	host        string
	port        string
	tlsCertFile string
	tlsKeyFile  string
}

// NewGRPCConfig creates new grpcConfig
func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	tlsServiceCertFile := os.Getenv(serviceTlsCertFile)
	if tlsServiceCertFile == "" {
		return nil, errors.New("grpc tls cert file not found")
	}

	tlsServiceKeyFile := os.Getenv(serviceTlsKeyFile)
	if tlsServiceKeyFile == "" {
		return nil, errors.New("grpc tls key file not found")
	}

	return &grpcConfig{
		host:        host,
		port:        port,
		tlsCertFile: tlsServiceCertFile,
		tlsKeyFile:  tlsServiceKeyFile,
	}, nil
}

// Address returns address from config
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// TlsCertFile returns path to TLS cert file from config
func (cfg *grpcConfig) TlsCertFile() string {
	return cfg.tlsCertFile
}

// TlsKeyFile returns path to TLS key file from config
func (cfg *grpcConfig) TlsKeyFile() string {
	return cfg.tlsKeyFile
}
