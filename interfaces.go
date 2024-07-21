package ax_tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	InitTracing(ctx context.Context, requestID, functionArn string) (func(context.Context) error, error)
	GetLogger() *slog.Logger
	StartSpan(ctx context.Context, name string) (context.Context, trace.Span)
	EndSpan(span trace.Span)
	AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue)
	LinkSpans(ctx context.Context, linkedCtx context.Context)
}

type TracerProvider interface {
	SetupTracer() (func(context.Context) error, error)
}
