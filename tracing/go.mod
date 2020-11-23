module github.com/HomesNZ/go-common/tracing

go 1.15

require (
	github.com/HomesNZ/go-common/env v0.0.0-20201120023436-5acd0c46a3d0
	github.com/pkg/errors v0.9.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
)
