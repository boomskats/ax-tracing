package ax_tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)


var (
    defaultTracer Tracer = NewDefaultTracer()
)

func InitTracing(ctx context.Context, requestID, functionArn string) (func(context.Context) error, error) {
    return defaultTracer.InitTracing(ctx, requestID, functionArn)
}

func GetLogger() *slog.Logger {
    return defaultTracer.GetLogger()
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
    return defaultTracer.StartSpan(ctx, name)
}

func EndSpan(span trace.Span) {
    defaultTracer.EndSpan(span)
}

func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
    defaultTracer.AddSpanEvent(ctx, name, attrs...)
}

func LinkSpans(ctx context.Context, linkedCtx context.Context) {
    defaultTracer.LinkSpans(ctx, linkedCtx)
}


