package ax_tracing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type envArg struct {
	name   string
	value  string
	defval string
}

type envTests []struct {
	name     string
	envs     envArg
	expected string
	assert   assert.ComparisonAssertionFunc
}

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
     

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            os.Setenv(tt.envs.name, tt.envs.value)
            actualValue := getEnv(tt.envs.name, tt.envs.defval)
            tt.assert(t, tt.expected, actualValue)
        })
    }
}
