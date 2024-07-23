package ax_tracing

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.opentelemetry.io/otel/attribute"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// TestIsTestMode tests the IsTestMode function
func TestIsTestMode(t *testing.T) {
    // Test with a regular context (should return false)
    assert.False(t, IsTestMode(context.Background()), "IsTestMode should return false for a regular context")

    // Test with a context that has test mode set (should return true)
    assert.True(t, IsTestMode(WithTestMode(context.Background())), "IsTestMode should return true for a context with test mode set")
}

func TestResource(t *testing.T) {
    res := Resource()
    assert.NotNil(t, res)
    attrs := res.Attributes()
    assert.Contains(t, attrs, semconv.ServiceNameKey.String(serviceName))
    assert.Contains(t, attrs, semconv.ServiceVersionKey.String(serviceVersion))
    assert.Contains(t, attrs, attribute.String("environment", deploymentEnvironment))
}

func TestInstallExportPipeline(t *testing.T) {
    ctx := context.Background()
    shutdown, err := InstallExportPipeline(ctx)
    assert.NoError(t, err)
    assert.NotNil(t, shutdown)

    err = shutdown(ctx)
    assert.NoError(t, err)
}

func TestWithTestMode(t *testing.T) {
    ctx := context.Background()
    testCtx := WithTestMode(ctx)
    assert.NotEqual(t, ctx, testCtx)
    assert.True(t, IsTestMode(testCtx))
    assert.False(t, IsTestMode(ctx))
}
