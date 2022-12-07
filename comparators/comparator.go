package comparators

import (
	"github.com/rini-labs/go-stream"
)

// Of creates a new Comparator[T] from the given function
func Of[T any](compare func(val1 T, val2 T) int) stream.Comparator[T] {
	return &comparator[T]{compare: compare}
}

type comparator[T any] struct {
	compare func(val1 T, val2 T) int
}

func (c *comparator[T]) Compare(val1 T, val2 T) int {
	return c.compare(val1, val2)
}

func NotOf[T any](c stream.Comparator[T]) stream.Comparator[T] {
	return &comparator[T]{
		compare: func(val1 T, val2 T) int {
			return -c.Compare(val1, val2)
		},
	}
}
