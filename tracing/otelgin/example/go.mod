module github.com/HomesNZ/go-common/tracing/otelgin/example

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/stdout v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
)
