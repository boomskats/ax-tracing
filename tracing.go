package ax_tracing

import (
	"context"
	"crypto/tls"
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

const (
)

var (
    serviceName  = os.Getenv("AXIOM_SERVICE_NAME")
	bearerToken  = os.Getenv("AXIOM_TOKEN")
	dataset      = os.Getenv("AXIOM_TRACES_DATASET") 
	otlpEndpoint = os.Getenv("AXIOM_URL")
	serviceVersion        = os.Getenv("AXIOM_SERVICE_VERSION")
	deploymentEnvironment = os.Getenv("AXIOM_ENVIRONMENT")
)

var tracer = otel.Tracer(serviceName)
var Logger = *slog.Default()

// InitTracing initializes both OpenTelemetry tracing and Axiom logging
func InitTracing(ctx context.Context, requestID, functionArn string) (func(context.Context) error, error) {
	// Set up Axiom logging
	lh, err := adapter.New()
	if err != nil {
		return nil, err
	}

	Logger := slog.New(lh).With("requestId", requestID).With("lambdaFunctionArn", functionArn)
	slog.SetDefault(Logger)
	Logger.Info("__ax-tracing logger initialised__")

	// Set up OpenTelemetry tracing
	otelShutdown, err := SetupTracer()
	if err != nil {
		Logger.Error("Failed to initialize OpenTelemetry", "error", err)
		return nil, err
	}
	Logger.Info("__ax-tracing otel tracer initialised__")

	// Return a combined shutdown function
	return func(shutdownCtx context.Context) error {
		if err := otelShutdown(shutdownCtx); err != nil {
			Logger.Error("Failed to shutdown OpenTelemetry", "error", err)
		}
		lh.Close()
		return nil
	}, nil
}

func GetLogger() *slog.Logger {
    return &Logger
}

// SetupTracer sets up the OpenTelemetry tracer
func SetupTracer() (func(context.Context) error, error) {
	ctx := context.Background()
	return InstallExportPipeline(ctx)
}

// Resource creates a new resource with service attributes
func Resource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
		attribute.String("environment", deploymentEnvironment),
	)
}

// InstallExportPipeline sets up the OpenTelemetry export pipeline
func InstallExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization":   bearerToken,
			"X-AXIOM-DATASET": dataset,
		}),
		otlptracehttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(Resource()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider.Shutdown, nil
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
