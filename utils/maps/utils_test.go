package maps

import (
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKeysFromMap(t *testing.T) {
	a := make(map[string][]byte)
	a["a"] = []byte("x")
	a["b"] = []byte("y")
	a["c"] = []byte("z")

	assert.True(t, strings.ArrayEqualElements([]string{"a", "b", "c"}, GetKeysFromMap(a)))
}

func TestGetValueFromMap(t *testing.T) {
	// Define test cases
	tests := []struct {
		name         string
		inputMap     map[string]int
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "Key exists in the map",
			inputMap:     map[string]int{"one": 1, "two": 2},
			key:          "one",
			defaultValue: 0,
			expected:     1,
		},
		{
			name:         "Key does not exist in the map",
			inputMap:     map[string]int{"one": 1, "two": 2},
			key:          "three",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "Empty map with a default value",
			inputMap:     map[string]int{},
			key:          "anykey",
			defaultValue: 99,
			expected:     99,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			got := GetValueFromMap(tc.inputMap, tc.key, tc.defaultValue)

			// Compare the result
			if got != tc.expected {
				t.Errorf("GetValueFromMap(%v, %q, %d) = %d; want %d",
					tc.inputMap, tc.key, tc.defaultValue, got, tc.expected)
			}
		})
	}
}
