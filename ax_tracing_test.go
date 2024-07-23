package ax_tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TestInitTracing(t *testing.T) {
	mockTracer := &MockTracer{}
	defaultTracer = mockTracer

	ctx := context.Background()
	requestID := "test-request-id"
	functionArn := "test-function-arn"

	mockShutdown := func(context.Context) error { return nil }
	mockTracer.On("InitTracing", ctx, requestID, functionArn).Return(mockShutdown, nil)

	shutdown, err := InitTracing(ctx, requestID, functionArn)
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)

	mockTracer.AssertExpectations(t)
}

func TestGetLogger(t *testing.T) {
	mockTracer := &MockTracer{}
	defaultTracer = mockTracer

	mockLogger := &slog.Logger{}
	mockTracer.On("GetLogger").Return(mockLogger)

	logger := GetLogger()
	assert.Equal(t, mockLogger, logger)

	mockTracer.AssertExpectations(t)
}

func TestStartSpan(t *testing.T) {
	mockTracer := &MockTracer{}
	defaultTracer = mockTracer

	ctx := context.Background()
	spanName := "test-span"
	mockSpan := &MockSpan{}
	mockTracer.On("StartSpan", ctx, spanName).Return(ctx, mockSpan)

	resultCtx, resultSpan := StartSpan(ctx, spanName)
	assert.Equal(t, ctx, resultCtx)
	assert.Equal(t, mockSpan, resultSpan)

	mockTracer.AssertExpectations(t)
}

func TestEndSpan(t *testing.T) {
	mockTracer := &MockTracer{}
	defaultTracer = mockTracer

	mockSpan := &MockSpan{}
	mockTracer.On("EndSpan", mockSpan).Return()

	EndSpan(mockSpan)

	mockTracer.AssertExpectations(t)
}

func TestAddSpanEvent(t *testing.T) {
	mockTracer := &MockTracer{}
	defaultTracer = mockTracer

	ctx := context.Background()
	eventName := "test-event"
	attrs := []attribute.KeyValue{attribute.String("key", "value")}

	mockTracer.On("AddSpanEvent", ctx, eventName, attrs).Return()

	AddSpanEvent(ctx, eventName, attrs...)

	mockTracer.AssertExpectations(t)
}

func TestLinkSpans(t *testing.T) {
	mockTracer := &MockTracer{}
	defaultTracer = mockTracer

	ctx1 := context.Background()
	ctx2 := context.Background()

	mockTracer.On("LinkSpans", ctx1, ctx2).Return()

	LinkSpans(ctx1, ctx2)

	mockTracer.AssertExpectations(t)
}
