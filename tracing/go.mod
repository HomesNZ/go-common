module github.com/HomesNZ/go-common/tracing

go 1.15

require (
	github.com/HomesNZ/go-common/env v0.0.0-20201123024206-a05abd62ee4e
	github.com/pkg/errors v0.9.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.14.0
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
)
