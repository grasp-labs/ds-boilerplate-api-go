package os

import "os"

// GetEnvOrDefault returns the value of an environment variable or a default value if the environment variable is not set
func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
