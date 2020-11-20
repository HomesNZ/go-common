package tracing

import (
	"context"
	"net/http"
	"net/url"

	"github.com/HomesNZ/go-common/env"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagators"
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
// this will also overwrite http.Default client with a tracing client, this is required to ensure that the trace context is passed onto dependencies
func InitTracer(ctx context.Context, cfg *TracerConfig, sampleType trace.Sampler) (func(), error) {
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

	global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))
	http.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithFilter(func(r *http.Request) bool {
				// ignore messages sent to the jaeger-agent
				return r.URL.Port() != "14268"
			}),
		),
	}
	http.DefaultTransport = http.DefaultClient.Transport

	tr := global.Tracer("init")
	_, span := tr.Start(ctx, "init")
	defer span.End()
	return flush, nil
}
