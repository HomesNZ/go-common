package tracing

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer adds a standard Jaeger Tracer to the otel global API
func InitTracer(ctx context.Context, serviceName, collectorHost, collectorPort, collectorPath string) (func(), error) {
	collectorEndpoint := fmt.Sprintf("http://%s:%s%s", collectorHost, collectorPort, collectorPath)
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(collectorEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: serviceName,
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
