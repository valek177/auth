package tracing

import (
	"github.com/uber/jaeger-client-go/config"

	internalCfg "github.com/valek177/auth/internal/config"
)

func Init(jaegerCfg internalCfg.JaegerConfig) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaegerCfg.SamplerType(),
			Param: jaegerCfg.SamplerParam(),
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: jaegerCfg.LocalAgentAddress(),
		},
	}

	_, err := cfg.InitGlobalTracer(jaegerCfg.ServiceName())
	if err != nil {
		return err
	}

	return nil
}
