package ax_tracing

import "context"

// DefaultTracerProvider is the default implementation of the TracerProvider interface.
type DefaultTracerProvider struct{}

var _ TracerProvider = (*DefaultTracerProvider)(nil) // Ensure DefaultTracerProvider implements TracerProvider interface

// SetupTracer initializes and sets up the OpenTelemetry tracer.
//
// It creates a background context and calls InstallExportPipeline to set up the tracer.
//
// Returns:
//   - A shutdown function that should be called to clean up resources when tracing is no longer needed
//   - An error if the setup fails
func (p *DefaultTracerProvider) SetupTracer() (func(context.Context) error, error) {
	ctx := context.Background()
	return InstallExportPipeline(ctx)
}
