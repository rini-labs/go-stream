package types

// Comparator function compares its two arguments for order. Returns a negative
// integer, zero, or a positive integer as the first argument is less than,
// equal to, or greater than the second.
type Comparator[T any] func(a, b T) int
