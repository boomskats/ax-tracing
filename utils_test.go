package ax_tracing

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestIsTestMode(t *testing.T) {
    assert.False(t, IsTestMode(context.Background()))
    assert.True(t, IsTestMode(WithTestMode(context.Background())))
}
