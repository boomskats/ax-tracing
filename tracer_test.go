package ax_tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestInitTracing(t *testing.T) {
    ctx := WithTestMode(context.Background())
    requestID := "test-request-id"
    functionArn := "test-function-arn"
    shutdownFunc, err := InitTracing(ctx, requestID, functionArn)
    assert.NoError(t, err)
    assert.NotNil(t, shutdownFunc)
}

func TestGetLogger(t *testing.T) {
    logger := GetLogger()
    assert.NotNil(t, logger)
}

func TestStartSpan(t *testing.T) {
    ctx := context.Background()
    name := "test-span"
    ctx, span := StartSpan(ctx, name)
    assert.NotNil(t, span)
    assert.NotNil(t, ctx)
}

func TestEndSpan(t *testing.T) {
    _, span := StartSpan(context.Background(), "test-span")
    span.End()
}

func TestAddSpanEvent(t *testing.T) {
    ctx := context.Background()
    _, span := StartSpan(ctx, "test-span")
    span.AddEvent("test-event")
}

func TestLinkSpans(t *testing.T) {
    ctx := context.Background()
    linkedCtx := context.Background()
    LinkSpans(ctx, linkedCtx)
}

func TestNewDefaultTracer(t *testing.T) {
    tracer := NewDefaultTracer()
    assert.NotNil(t, tracer.tracerProvider)
    assert.NotNil(t, tracer.tracer)
    assert.NotNil(t, tracer.logger)
}
