package ax_tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
    ctx := WithTestMode(context.Background())
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
    ctx := WithTestMode(context.Background())
    _, span := StartSpan(ctx, "test-span")
    span.AddEvent("test-event")
}

func TestLinkSpans(t *testing.T) {
    // Create a mock tracer
    mockTracer := &MockTracer{}

    // Set up expectations
    mockTracer.On("LinkSpans", mock.Anything, mock.Anything).Return()

    // Create two contexts with spans
    ctx1 := WithTestMode(context.Background())
    ctx2 := WithTestMode(context.Background())

    // Call LinkSpans with the mock tracer
    mockTracer.LinkSpans(ctx1, ctx2)

    // Assert that LinkSpans was called with the correct parameters
    mockTracer.AssertCalled(t, "LinkSpans", ctx1, ctx2)

    // Verify all expectations were met
    mockTracer.AssertExpectations(t)
}

func TestNewDefaultTracer(t *testing.T) {
    tracer := NewDefaultTracer()
    assert.NotNil(t, tracer.tracerProvider)
    assert.NotNil(t, tracer.tracer)
    assert.NotNil(t, tracer.logger)
}

func TestDefaultTracer_InitTracing(t *testing.T) {
	ctx := WithTestMode(context.Background())
	requestID := "test-request-id"
	functionArn := "test-function-arn"

	tracer := NewDefaultTracer()

	shutdown, err := tracer.InitTracing(ctx, requestID, functionArn)
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)

	err = shutdown(ctx)
	assert.NoError(t, err)
}

func TestDefaultTracer_GetLogger(t *testing.T) {
	tracer := NewDefaultTracer()
	logger := tracer.GetLogger()
	assert.NotNil(t, logger)
}

func TestDefaultTracer_StartSpan(t *testing.T) {
	tracer := NewDefaultTracer()
	ctx := WithTestMode(context.Background())
	spanName := "test-span"

	newCtx, span := tracer.StartSpan(ctx, spanName)
	assert.NotNil(t, newCtx)
	assert.NotNil(t, span)

	tracer.EndSpan(span)
}

func TestDefaultTracer_AddSpanEvent(t *testing.T) {
	tracer := NewDefaultTracer()
	ctx := WithTestMode(context.Background())
	eventName := "test-event"
	attr := attribute.String("key", "value")

	_, span := tracer.StartSpan(ctx, "parent-span")
	ctx = trace.ContextWithSpan(ctx, span)

	tracer.AddSpanEvent(ctx, eventName, attr)

	tracer.EndSpan(span)
}

func TestDefaultTracer_LinkSpans(t *testing.T) {
	tracer := NewDefaultTracer()
	ctx1 := WithTestMode(context.Background())
	ctx2 := WithTestMode(context.Background())

	_, span1 := tracer.StartSpan(ctx1, "span1")
	ctx1 = trace.ContextWithSpan(ctx1, span1)

	_, span2 := tracer.StartSpan(ctx2, "span2")
	ctx2 = trace.ContextWithSpan(ctx2, span2)

	tracer.LinkSpans(ctx1, ctx2)

	tracer.EndSpan(span1)
	tracer.EndSpan(span2)
}
