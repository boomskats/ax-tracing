package ax_tracing

import (
	"context"
	"log/slog"

	adapter "github.com/axiomhq/axiom-go/adapters/slog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DefaultTracer struct {
	logger         *slog.Logger
	tracer         trace.Tracer
	tracerProvider TracerProvider
}

// InitTracing initializes OpenTelemetry tracing and Axiom logging
// and returns a shutdown function which should at least be deferred.

func (t *DefaultTracer) InitTracing(
	ctx context.Context,
	requestID, functionArn string,
) (
	func(context.Context) error,
	error,
) {
    // If this is running in test mode, then create a no-op logger
    // and return a no-op shutdown function. 
    if IsTestMode(ctx) {
        t.logger = slog.Default()
        t.logger.Info("__ax-tracing initialised in test mode__")
        return func(context.Context) error {
            t.logger.Info("__ax-tracing test mode shutdown__")
            return nil
        }, nil
    }

	// Set up Axiom logging
	lh, err := adapter.New()
	if err != nil {
		return nil, err
	}

	t.logger = slog.New(lh).With(
        "requestId", requestID).With(
        "lambdaFunctionArn", functionArn)
	slog.SetDefault(t.logger)
	t.logger.Info("__ax-tracing logger initialised__")

	// Set up OpenTelemetry tracing
	otelShutdown, err := t.tracerProvider.SetupTracer()
	if err != nil {
		t.logger.Error(
            "Failed to initialize OpenTelemetry", 
            "error", err)
		return nil, err
	}
	t.logger.Info("__ax-tracing otel tracer initialised__")

	// Return a combined shutdown function
	return func(shutdownCtx context.Context) error {
		if err := otelShutdown(shutdownCtx); err != nil {
			t.logger.Error(
                "Failed to shutdown OpenTelemetry", 
                "error", err)
		}
		t.logger.Info("__ax-tracing__ shutdown complete")
		lh.Close()
		return nil
	}, nil
}

// GetLogger returns the slog default logger. It
// is probably not that useful given how the
// logger is globally initialised
func (t *DefaultTracer) GetLogger() *slog.Logger {
	return t.logger
}

// StartSpan starts a new span and returns the context and span
func (t *DefaultTracer) StartSpan(
	ctx context.Context, name string,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name)
}

// EndSpan ends the given span
func (t *DefaultTracer) EndSpan(span trace.Span) {
	span.End()
}

// AddSpanEvent adds an event to the current span and logs
// the fact that it has added it. This is useful for
// debugging but not much else.
func (t *DefaultTracer) AddSpanEvent(
    ctx context.Context, 
    name string, 
    attrs ...attribute.KeyValue,
) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
	t.logger.InfoContext(
        ctx, "Span event added", "event", 
        name, "attributes", attrs)
}

// LinkSpans creates a link between current span and provided context
func (t *DefaultTracer) LinkSpans(
    ctx context.Context, 
    linkedCtx context.Context,
) {
	span := trace.SpanFromContext(ctx)
	linkedSpan := trace.SpanFromContext(linkedCtx)
	span.AddLink(trace.Link{SpanContext: linkedSpan.SpanContext()})
	t.logger.InfoContext(ctx, "Spans linked")
}

// NewDefaultTracer creates a new default tracer
func NewDefaultTracer() *DefaultTracer {
	return &DefaultTracer{
		logger:         slog.Default(),
		tracer:         otel.Tracer(serviceName),
		tracerProvider: &DefaultTracerProvider{},
	}
}

