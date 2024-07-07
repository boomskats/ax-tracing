// Package ax_tracing provies a simplified otel
// initialization for Lambda functions, integrating
// with Axiom for logging and tracing.
package ax_tracing

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	adapter "github.com/axiomhq/axiom-go/adapters/slog"
)

type Config struct {
	ServiceName           string
	BearerToken           string
	Dataset               string
	OTLPEndpoint          string
	ServiceVersion        string
	DeploymentEnvironment string
}

var config Config

var (
	ErrInitLogger = errors.New("failed to initialize logger")
	ErrInitTracer = errors.New("failed to initialize OpenTelemetry tracer")

	tracer trace.Tracer
	Logger = *slog.Default()
)

func init() {
	config = Config{
		ServiceName:           os.Getenv("AXIOM_SERVICE_NAME"),
		BearerToken:           "Bearer " + os.Getenv("AXIOM_TOKEN"),
		Dataset:               os.Getenv("AXIOM_TRACES_DATASET"),
		OTLPEndpoint:          os.Getenv("AXIOM_OTLP_ENDPOINT"),
		ServiceVersion:        os.Getenv("AXIOM_SERVICE_VERSION"),
		DeploymentEnvironment: os.Getenv("AXIOM_ENVIRONMENT"),
	}
}

// InitTracing initializes OpenTelemetry tracer and
// Axiom logging slog logger, and returns a combined
// shutdown function and an error if init fails
//
// The shutdown function should be deferred in the
// main function to ensure proper cleanup of resources
func InitTracing(
	ctx context.Context,
	requestID, functionArn string,
) (func(context.Context) error, error) {
	// Set up Axiom logging
	lh, err := adapter.New()
	if err != nil {
		return nil, err
	}

	logger := slog.New(lh).With("requestId", requestID).With("lambdaFunctionArn", functionArn)
	slog.SetDefault(logger)
	slog.Debug("__ax-tracing logger initialised__")

	// Set up OpenTelemetry tracing
	otelShutdown, err := installExportPipeline(ctx)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
		return nil, err
	}
	Logger.Debug("__ax-tracing otel tracer initialised__")

	// Return a combined shutdown function
	return func(shutdownCtx context.Context) error {
		if err := otelShutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown OpenTelemetry", "error", err)
		}
		lh.Close()
		slog.Debug("__ax-tracing logger shutdown__")
		return nil
	}, nil
}

// GetLogger returns the slog default logger.
// It is probably not that useful given how 
// the logger is initialised
func GetLogger() *slog.Logger {
	return &Logger
}

// createResource creates a new resource with service attributes
func createResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
		attribute.String("environment", config.DeploymentEnvironment),
	)
}

// installExportPipeline sets up the OpenTelemetry export pipeline
func installExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.OTLPEndpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization":   config.BearerToken,
			"X-AXIOM-DATASET": config.Dataset,
		}),
		otlptracehttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(createResource()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Initialize the tracer after setting up the tracer provider
	tracer = otel.GetTracerProvider().Tracer(config.ServiceName)

	return func(shutdownCtx context.Context) error {
		// Force flush any remaining spans
		if err := tracerProvider.ForceFlush(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown OpenTelemetry", "error", err)
			return err
		}
		slog.Error("__otel tracer flushed and shutdown__")
		return tracerProvider.Shutdown(shutdownCtx)
	}, nil
}

// StartSpan starts a new span and returns the context and span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

// EndSpan ends the given span
func EndSpan(span trace.Span) {
	span.End()
}

// TracedFunction is a wrapper for functions that need tracing
func TracedFunction(ctx context.Context, name string, f func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer EndSpan(span)

	err := f(ctx)
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Error in traced function", "function", name, "error", err)
	}

	return err
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
	slog.InfoContext(ctx, "Span event added", "event", name, "attributes", attrs)
}

// SetSpanAttributes sets attributes on the current span
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
	slog.InfoContext(ctx, "Span attributes set", "attributes", attrs)
}

// LinkSpans creates a link between the current span and the provided context
func LinkSpans(ctx context.Context, linkedCtx context.Context) {
	span := trace.SpanFromContext(ctx)
	linkedSpan := trace.SpanFromContext(linkedCtx)
	span.AddLink(trace.Link{SpanContext: linkedSpan.SpanContext()})
	slog.InfoContext(ctx, "Spans linked")
}
