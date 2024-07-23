package ax_tracing

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

// TestIsTestMode tests the IsTestMode function
func TestIsTestMode(t *testing.T) {
    // Test with a regular context (should return false)
    assert.False(t, IsTestMode(context.Background()), "IsTestMode should return false for a regular context")

    // Test with a context that has test mode set (should return true)
    assert.True(t, IsTestMode(WithTestMode(context.Background())), "IsTestMode should return true for a context with test mode set")
}
