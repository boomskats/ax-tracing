package ax_tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTracerProvider_SetupTracer(t *testing.T) {
	provider := &DefaultTracerProvider{}

	shutdown, err := provider.SetupTracer()
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)

	err = shutdown(context.Background())
	assert.NoError(t, err)
}
