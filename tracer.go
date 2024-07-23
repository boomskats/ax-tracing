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

var _ Tracer = (*DefaultTracer)(nil) // Ensure DefaultTracer implements Tracer interface

// InitTracing initializes OpenTelemetry tracing and Axiom logging.
// It sets up the logger with the provided requestID and functionArn,
// and initializes the OpenTelemetry tracer.
//
// Parameters:
//   - ctx: The context for the operation
//   - requestID: A unique identifier for the request
//   - functionArn: The ARN of the Lambda function
//
// Returns:
//   - A shutdown function that should be deferred to clean up resources
//   - An error if initialization fails
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

// GetLogger returns the slog default logger.
// Note: This may not be very useful as the logger is globally initialized.
//
// Returns:
//   - The default slog.Logger
func (t *DefaultTracer) GetLogger() *slog.Logger {
	return t.logger
}

// StartSpan starts a new span and returns the updated context and the span.
//
// Parameters:
//   - ctx: The parent context
//   - name: The name of the span
//
// Returns:
//   - The updated context containing the new span
//   - The newly created span
func (t *DefaultTracer) StartSpan(
	ctx context.Context, name string,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name)
}

// EndSpan ends the given span.
//
// Parameters:
//   - span: The span to end
func (t *DefaultTracer) EndSpan(span trace.Span) {
	span.End()
}

// AddSpanEvent adds an event to the current span and logs the fact that it has been added.
// This can be useful for debugging purposes.
//
// Parameters:
//   - ctx: The context containing the current span
//   - name: The name of the event
//   - attrs: Optional attributes to add to the event
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

// LinkSpans creates a link between the current span and the span in the provided context.
//
// Parameters:
//   - ctx: The context containing the current span
//   - linkedCtx: The context containing the span to link to
func (t *DefaultTracer) LinkSpans(
    ctx context.Context, 
    linkedCtx context.Context,
) {
	span := trace.SpanFromContext(ctx)
	linkedSpan := trace.SpanFromContext(linkedCtx)
	span.AddLink(trace.Link{SpanContext: linkedSpan.SpanContext()})
	t.logger.InfoContext(ctx, "Spans linked")
}

// NewDefaultTracer creates and returns a new instance of DefaultTracer.
//
// Returns:
//   - A pointer to a new DefaultTracer instance
func NewDefaultTracer() *DefaultTracer {
	return &DefaultTracer{
		logger:         slog.Default(),
		tracer:         otel.Tracer(serviceName),
		tracerProvider: &DefaultTracerProvider{},
	}
}

