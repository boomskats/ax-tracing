package ax_tracing

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is an interface that defines the methods for tracing operations.
type Tracer interface {
	// InitTracing initializes the tracing system with the given request ID and function ARN.
	// It returns a shutdown function and an error if initialization fails.
	InitTracing(ctx context.Context, requestID, functionArn string) (func(context.Context) error, error)

	// GetLogger returns the logger instance used by the tracer.
	GetLogger() *slog.Logger

	// StartSpan starts a new span with the given name and returns the updated context and the span.
	StartSpan(ctx context.Context, name string) (context.Context, trace.Span)

	// EndSpan ends the given span.
	EndSpan(span trace.Span)

	// AddSpanEvent adds an event to the current span in the given context.
	AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue)

	// LinkSpans creates a link between the current span and the span in the provided context.
	LinkSpans(ctx context.Context, linkedCtx context.Context)
}

// TracerProvider is an interface for creating and setting up tracers.
type TracerProvider interface {
	// SetupTracer initializes and sets up the tracer.
	// It returns a shutdown function and an error if setup fails.
	SetupTracer() (func(context.Context) error, error)
}
