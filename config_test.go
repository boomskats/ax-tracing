package ax_tracing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		envName  string
		envValue string
		defValue string
		expected string
	}{
		{"AXIOM_SERVICE_NAME", "AXIOM_SERVICE_NAME", "test-service", "default-ax-service", "test-service"},
		{"AXIOM_SERVICE_VERSION", "AXIOM_SERVICE_VERSION", "1.0.0", "0.0.0", "1.0.0"},
		{"AXIOM_ENVIRONMENT", "AXIOM_ENVIRONMENT", "test-env", "default-ax-environment", "test-env"},
		{"AXIOM_TOKEN", "AXIOM_TOKEN", "test-token", "default-ax-token", "test-token"},
		{"AXIOM_TRACES_DATASET", "AXIOM_TRACES_DATASET", "test-dataset", "default-ax-dataset", "test-dataset"},
		{"AXIOM_OTLP_ENDPOINT", "AXIOM_OTLP_ENDPOINT", "test-otlp-endpoint", "default-ax-otlp-endpoint", "test-otlp-endpoint"},
		{"NON_EXISTENT_VAR", "NON_EXISTENT_VAR", "", "fallback", "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envName, tt.envValue)
				defer os.Unsetenv(tt.envName)
			}
			actual := getEnv(tt.envName, tt.defValue)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestConfigVariables(t *testing.T) {
	// Test default values
	assert.Equal(t, "default-ax-service", serviceName)
	assert.Equal(t, "Bearer ", bearerToken)
	assert.Equal(t, "", dataset)
	assert.Equal(t, "", otlpEndpoint)
	assert.Equal(t, "0.0.0", serviceVersion)
	assert.Equal(t, "default-ax-environment", deploymentEnvironment)

	// Set environment variables
	envVars := map[string]string{
		"AXIOM_SERVICE_NAME":    "test-service",
		"AXIOM_TOKEN":           "test-token",
		"AXIOM_TRACES_DATASET":  "test-dataset",
		"AXIOM_OTLP_ENDPOINT":   "test-endpoint",
		"AXIOM_SERVICE_VERSION": "1.0.0",
		"AXIOM_ENVIRONMENT":     "test-environment",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	// Reinitialize variables
	serviceName = getEnv("AXIOM_SERVICE_NAME", "default-ax-service")
	bearerToken = "Bearer " + getEnv("AXIOM_TOKEN", "")
	dataset = getEnv("AXIOM_TRACES_DATASET", "")
	otlpEndpoint = getEnv("AXIOM_OTLP_ENDPOINT", "")
	serviceVersion = getEnv("AXIOM_SERVICE_VERSION", "0.0.0")
	deploymentEnvironment = getEnv("AXIOM_ENVIRONMENT", "default-ax-environment")

	// Test new values
	assert.Equal(t, "test-service", serviceName)
	assert.Equal(t, "Bearer test-token", bearerToken)
	assert.Equal(t, "test-dataset", dataset)
	assert.Equal(t, "test-endpoint", otlpEndpoint)
	assert.Equal(t, "1.0.0", serviceVersion)
	assert.Equal(t, "test-environment", deploymentEnvironment)
}
