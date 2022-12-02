package comparators

import (
	"github.com/rini-labs/go-stream/pkg/types"
	"strings"

	"golang.org/x/exp/constraints"
)

// Natural implements the Comparator for those elements whose type
// has a natural order (numbers and strings)
func Natural[T constraints.Ordered](a, b T) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return +1
}

// Int implements the Comparator for signed integers. This will be usually
// faster than Natural comparator
func Int[T constraints.Integer](a, b T) int {
	return int(a - b)
}

// IgnoreCase implements order.Comparator for strings, ignoring the case.
func IgnoreCase(a, b string) int {
	return Natural(strings.ToLower(a), strings.ToLower(b))
}

// Inverse result of the Comparator function for inverted sorts
func Inverse[T any](cmp types.Comparator[T]) types.Comparator[T] {
	return func(a, b T) int {
		return -cmp(a, b)
	}
}
