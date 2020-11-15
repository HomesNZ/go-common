package tracing

import (
	"context"
	"fmt"
	"time"

	"github.com/HomesNZ/env"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TracerConfig manages initialization of the tracing collector
type TracerConfig struct {
	Name          string
	CollectorHost string
	CollectorPort string
	CollectorPath string
}

func (cfg TracerConfig) collectorEndpoint() string {
	return fmt.Sprintf("http://%s:%s%s", cfg.CollectorHost, cfg.CollectorPort, cfg.CollectorPath)
}

// TracerConfigFromEnv initializes a TracerConfig from environment variables.
func TracerConfigFromEnv() TracerConfig {
	return TracerConfig{
		Name:          env.MustGetString("SERVICE_NAME"),
		CollectorHost: env.MustGetString("TRACING_COLLECTOR_HOST"),
		CollectorPort: env.MustGetString("TRACING_COLLECTOR_PORT"),
		CollectorPath: env.MustGetString("TRACING_COLLECTOR_PATH"),
	}
}

// InitTracer adds a standard Jaeger Tracer to the otel global API
func InitTracer(ctx context.Context, cfg TracerConfig) (func(), error) {
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.collectorEndpoint()),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: cfg.Name,
		}),
		jaeger.WithSDK(&trace.Config{DefaultSampler: trace.AlwaysSample()}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "InstallNewPipeline")
	}

	tr := global.Tracer("init")
	_, span := tr.Start(ctx, "init")
	defer span.End()
	time.Sleep(1 * time.Second)
	return flush, nil
}
