package comparators

import (
	"strings"

	"github.com/rini-labs/go-stream"
	"golang.org/x/exp/constraints"
)

// Natural implements the Comparator for those elements whose type
// has a natural order (numbers and strings)
func Natural[T constraints.Ordered]() stream.Comparator[T] {
	return Of[T](func(a, b T) int {
		if a == b {
			return 0
		}
		if a < b {
			return -1
		}
		return +1
	})
}

// Int implements the Comparator for signed integers. This will be usually
// faster than Natural comparator
func Int[T constraints.Integer]() stream.Comparator[T] {
	return Of[T](func(a, b T) int {
		return int(a - b)
	})
}

// IgnoreCase implements order.Comparator for strings, ignoring the case.
func IgnoreCase() stream.Comparator[string] {
	return Of[string](func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})
}
