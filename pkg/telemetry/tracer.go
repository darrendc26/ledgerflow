package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(serviceName string) (*trace.TracerProvider, error) {
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
