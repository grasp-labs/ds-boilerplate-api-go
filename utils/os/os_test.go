package os

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	// Set up test cases
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			envKey:       "TEST_KEY_EXISTS",
			envValue:     "value_from_env",
			defaultValue: "default_value",
			expected:     "value_from_env",
		},
		{
			name:         "Environment variable does not exist",
			envKey:       "TEST_KEY_MISSING",
			envValue:     "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set the environment variable if a value is provided
			if tc.envValue != "" {
				err := os.Setenv(tc.envKey, tc.envValue)
				if err != nil {
					t.Fatalf("Failed to set environment variable: %v", err)
				}
				// Clean up after the test
				defer os.Unsetenv(tc.envKey)
			} else {
				// Ensure the environment variable is not set
				os.Unsetenv(tc.envKey)
			}

			// Call the function
			got := GetEnvOrDefault(tc.envKey, tc.defaultValue)

			// Compare the result with the expected value
			if got != tc.expected {
				t.Errorf("GetEnvOrDefault(%q, %q) = %q; want %q", tc.envKey, tc.defaultValue, got, tc.expected)
			}
		})
	}
}
