package ax_tracing

import "context"

type DefaultTracerProvider struct{}

var _ TracerProvider = (*DefaultTracerProvider)(nil) // Ensure DefaultTracerProvider implements TracerProvider interface

func (p *DefaultTracerProvider) SetupTracer() (func(context.Context) error, error) {
	ctx := context.Background()
	return InstallExportPipeline(ctx)
}
