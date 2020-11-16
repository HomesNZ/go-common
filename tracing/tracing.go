package tracing

import (
	"context"
	"net/url"
	"time"

	"github.com/HomesNZ/env"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TracerConfig manages initialization of the tracing collector
type TracerConfig struct {
	Name              string
	CollectorEndpoint string
}

// TracerConfigFromEnv initializes a TracerConfig from environment variables.
func TracerConfigFromEnv() (*TracerConfig, error) {
	endpoint := jaeger.CollectorEndpointFromEnv()
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Parse failed")
	}
	return &TracerConfig{
		Name:              env.MustGetString("SERVICE_NAME"),
		CollectorEndpoint: endpoint,
	}, nil
}

// InitTracer adds a standard Jaeger Tracer to the otel global API
func InitTracer(ctx context.Context, cfg TracerConfig, sampleType trace.Sampler) (func(), error) {
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.CollectorEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: cfg.Name,
		}),
		jaeger.WithSDK(&trace.Config{DefaultSampler: sampleType}),
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
