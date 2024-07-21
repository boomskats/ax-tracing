package ax_tracing

import "context"

type DefaultTracerProvider struct{}

func (p *DefaultTracerProvider) SetupTracer() (func(context.Context) error, error) {
	ctx := context.Background()
	return InstallExportPipeline(ctx)
}

