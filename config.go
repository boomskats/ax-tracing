package ax_tracing

import "os"

var (
	serviceName           = getEnv("AXIOM_SERVICE_NAME", "default-ax-service")
	bearerToken           = "Bearer " + getEnv("AXIOM_TOKEN", "")
	dataset               = getEnv("AXIOM_TRACES_DATASET", "")
	otlpEndpoint          = getEnv("AXIOM_OTLP_ENDPOINT", "")
	serviceVersion        = getEnv("AXIOM_SERVICE_VERSION", "0.0.0")
	deploymentEnvironment = getEnv("AXIOM_ENVIRONMENT", "default-ax-environment")
)


func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
