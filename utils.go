package ax_tracing

import (
	"context"
	"crypto/tls"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Resource creates and returns a new resource with service and environment attributes.
//
// Returns:
//   - A pointer to a resource.Resource with predefined attributes
func Resource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
		attribute.String("environment", deploymentEnvironment),
	)
}

// InstallExportPipeline sets up the OpenTelemetry export pipeline.
//
// It configures the OTLP HTTP exporter, creates a new TracerProvider with the exporter,
// and sets up the global OpenTelemetry instance.
//
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - A shutdown function that should be called to clean up resources when tracing is no longer needed
//   - An error if the setup fails
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

	return func(shutdownCtx context.Context) error {
		// Force flush any remaining spans
		if err := tracerProvider.ForceFlush(shutdownCtx); err != nil {
			slog.Error("Failed to force flush spans", "error", err)
			return err
		}
		slog.Info("Spans flushed")
		return tracerProvider.Shutdown(shutdownCtx)
	}, nil
}

type testModeKey struct{}

// WithTestMode returns a new context with a test mode flag set.
//
// This function is used to create a context that indicates the code is running in test mode.
//
// Parameters:
//   - ctx: The parent context
//
// Returns:
//   - A new context with the test mode flag set
func WithTestMode(ctx context.Context) context.Context {
    return context.WithValue(ctx, testModeKey{}, true)
}

// IsTestMode checks if the given context has the test mode flag set.
//
// Parameters:
//   - ctx: The context to check
//
// Returns:
//   - true if the context has the test mode flag set, false otherwise
func IsTestMode(ctx context.Context) bool {
    return ctx.Value(testModeKey{}) != nil
}
