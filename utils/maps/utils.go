package maps

import (
	"fmt"
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/slice"
	"strings"
)

// GetKeysFromMap returns the keys of a map
func GetKeysFromMap[T comparable, V any](m map[T]V) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MergeMaps combines two maps into one map.
// Maps are merged from the beginning of the sources slice,
// in case of key conflicts, values from source map with higher index will overwrite the previous ones.
func MergeMaps[T comparable, V any](sources ...map[T]V) map[T]V {
	result := make(map[T]V)
	for _, source := range sources {
		for k, v := range source {
			result[k] = v
		}
	}
	return result
}

// FromString Converts a string with comma separated key-value pairs to a map[string]string
func FromString(keyValuePairs string) map[string]string {
	return slice.Reduce(strings.Split(keyValuePairs, ","), make(map[string]string), func(acc map[string]string, pair string) map[string]string {
		keyValue := strings.Split(pair, "=")
		if len(keyValue) != 2 {
			return acc
		}
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		if len(key) > 0 && len(value) > 0 {
			acc[key] = value
		}
		return acc
	})
}

// ToString Converts map[string]string to a string with comma separated key-vaults, spaces are trimmed
func ToString(keyValues map[string]string) string {
	return strings.Join(slice.Reduce(GetKeysFromMap(keyValues), make([]string, 0), func(acc []string, key string) []string {
		trimmedKey := strings.TrimSpace(key)
		if len(trimmedKey) == 0 {
			return acc
		}
		return append(acc, fmt.Sprintf("%s=%s", trimmedKey, strings.TrimSpace(keyValues[key])))
	}), ",")
}

// GetValueFromMap Returns the value of a key from a map If the key does not
// exist, the default value is returned
func GetValueFromMap[T comparable, V any](m map[T]V, key T, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}
