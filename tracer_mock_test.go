package ax_tracing

import (
	"context"
	"log/slog"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// MockTracer is a mock implementation of the Tracer interface.
// It can be used by both package consumers and package developers
// to test code that depends on the Tracer interface without setting up
// a real tracing infrastructure.
//
// Usage example:
//
//	func TestSomeFunctionUsingTracer(t *testing.T) {
//		mockTracer := &MockTracer{}
//		mockTracer.On("StartSpan", mock.Anything, "span-name").Return(context.Background(), &MockSpan{})
//		
//		// Use mockTracer in your test...
//		
//		mockTracer.AssertExpectations(t)
//	}
type MockTracer struct {
	mock.Mock
}

func (m *MockTracer) InitTracing(ctx context.Context, requestID, functionArn string) (func(context.Context) error, error) {
	args := m.Called(ctx, requestID, functionArn)
	return args.Get(0).(func(context.Context) error), args.Error(1)
}

func (m *MockTracer) GetLogger() *slog.Logger {
	args := m.Called()
	return args.Get(0).(*slog.Logger)
}

func (m *MockTracer) StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	args := m.Called(ctx, name)
	return args.Get(0).(context.Context), args.Get(1).(trace.Span)
}

func (m *MockTracer) EndSpan(span trace.Span) {
	m.Called(span)
}

func (m *MockTracer) AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	m.Called(ctx, name, attrs)
}

func (m *MockTracer) LinkSpans(ctx context.Context, linkedCtx context.Context) {
	m.Called(ctx, linkedCtx)
}
