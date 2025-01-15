package slice

// Reduce applies a function against an accumulator and each element in the slice (from left to right) to reduce it to a single value.
func Reduce[TSource, TAccumulation any](source []TSource, seed TAccumulation, accumulator func(TAccumulation, TSource) TAccumulation) TAccumulation {
	for _, value := range source {
		seed = accumulator(seed, value)
	}
	return seed
}

// Any checks if any element in the slice satisfies the condition
func Any[T any](source []T, condition func(T) bool) bool {
	for _, value := range source {
		if condition(value) {
			return true
		}
	}
	return false
}

// All checks if all elements in the slice satisfy the condition
func All[T any](source []T, condition func(T) bool) bool {
	return !Any(source, func(v T) bool { return !condition(v) })
}

// FindAll returns all elements in the slice that satisfy the condition
func FindAll[T any](source []T, condition func(T) bool) []T {
	result := make([]T, 0)
	for _, value := range source {
		if condition(value) {
			result = append(result, value)
		}
	}
	return result
}

// FindFirst returns the first element in the slice that satisfies the condition
func FindFirst[T any](source []T, condition func(T) bool) (element T, ok bool) {
	for _, value := range source {
		if condition(value) {
			element = value
			ok = true
			return
		}
	}
	return
}
