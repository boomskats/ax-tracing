package ax_tracing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// envArg represents an environment variable with its name, value, and default value
type envArg struct {
	name   string
	value  string
	defval string
}

// envTests is a slice of test cases for environment variables
type envTests []struct {
	name     string
	envs     envArg
	expected string
	assert   assert.ComparisonAssertionFunc
}

// TestGetEnv tests the getEnv function with various environment variables
func TestGetEnv(t *testing.T) {
    tests := envTests{
        {
            name: "test-1",
            envs: envArg{
                name:   "AXIOM_SERVICE_NAME",
                value:  "test-service",
                defval: "default-ax-service",
            },
            expected: "test-service",
            assert:   assert.Equal,
        },
        {
            name: "test-2",
            envs: envArg{
                name:   "AXIOM_SERVICE_VERSION",
                value:  "1.0.0",
                defval: "0.0.0",
            },
            expected: "1.0.0",
            assert:   assert.Equal,
        },
        {
            name: "test-3",
            envs: envArg{
                name:   "AXIOM_ENVIRONMENT",
                value:  "test-env",
                defval: "default-ax-environment",
            },
            expected: "test-env",
            assert:   assert.Equal,
        },
        {
            name: "test-4",
            envs: envArg{
                name:   "AXIOM_TOKEN",
                value:  "test-token",
                defval: "default-ax-token",
            },
            expected: "test-token",
            assert:   assert.Equal,
        },
        {
            name: "test-5",
            envs: envArg{
                name:   "AXIOM_TRACES_DATASET",
                value:  "test-dataset",
                defval: "default-ax-dataset",
            },
            expected: "test-dataset",
            assert:   assert.Equal,
        },
        {
            name: "test-6",
            envs: envArg{
                name:   "AXIOM_OTLP_ENDPOINT",
                value:  "test-otlp-endpoint",
                defval: "default-ax-otlp-endpoint",
            },
            expected: "test-otlp-endpoint",
            assert:   assert.Equal,
        },
    }
     
    // Iterate through all test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set the environment variable for the test
            os.Setenv(tt.envs.name, tt.envs.value)
            // Get the value using getEnv
            actualValue := getEnv(tt.envs.name, tt.envs.defval)
            // Assert that the actual value matches the expected value
            tt.assert(t, tt.expected, actualValue)
        })
    }
}
package ax_tracing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	// Test when environment variable is set
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	value := getEnv("TEST_ENV_VAR", "fallback")
	assert.Equal(t, "test_value", value)

	// Test when environment variable is not set
	value = getEnv("NON_EXISTENT_VAR", "fallback")
	assert.Equal(t, "fallback", value)
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
	os.Setenv("AXIOM_SERVICE_NAME", "test-service")
	os.Setenv("AXIOM_TOKEN", "test-token")
	os.Setenv("AXIOM_TRACES_DATASET", "test-dataset")
	os.Setenv("AXIOM_OTLP_ENDPOINT", "test-endpoint")
	os.Setenv("AXIOM_SERVICE_VERSION", "1.0.0")
	os.Setenv("AXIOM_ENVIRONMENT", "test-environment")
	defer func() {
		os.Unsetenv("AXIOM_SERVICE_NAME")
		os.Unsetenv("AXIOM_TOKEN")
		os.Unsetenv("AXIOM_TRACES_DATASET")
		os.Unsetenv("AXIOM_OTLP_ENDPOINT")
		os.Unsetenv("AXIOM_SERVICE_VERSION")
		os.Unsetenv("AXIOM_ENVIRONMENT")
	}()

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
