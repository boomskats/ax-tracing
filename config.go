package ax_tracing

import "os"

// Configuration variables for the ax_tracing package
var (
	// serviceName is the name of the service, retrieved from the AXIOM_SERVICE_NAME environment variable
	serviceName = getEnv("AXIOM_SERVICE_NAME", "default-ax-service")

	// bearerToken is the authentication token for Axiom, retrieved from the AXIOM_TOKEN environment variable
	bearerToken = "Bearer " + getEnv("AXIOM_TOKEN", "")

	// dataset is the name of the Axiom dataset to use, retrieved from the AXIOM_TRACES_DATASET environment variable
	dataset = getEnv("AXIOM_TRACES_DATASET", "")

	// otlpEndpoint is the endpoint for the OpenTelemetry collector, retrieved from the AXIOM_OTLP_ENDPOINT environment variable
	otlpEndpoint = getEnv("AXIOM_OTLP_ENDPOINT", "")

	// serviceVersion is the version of the service, retrieved from the AXIOM_SERVICE_VERSION environment variable
	serviceVersion = getEnv("AXIOM_SERVICE_VERSION", "0.0.0")

	// deploymentEnvironment is the deployment environment, retrieved from the AXIOM_ENVIRONMENT environment variable
	deploymentEnvironment = getEnv("AXIOM_ENVIRONMENT", "default-ax-environment")
)

// getEnv retrieves the value of an environment variable, or returns a fallback value if it's not set.
//
// Parameters:
//   - key: The name of the environment variable
//   - fallback: The default value to return if the environment variable is not set
//
// Returns:
//   - The value of the environment variable, or the fallback value if not set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
