package iterator

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/pkg/types"
)

func Empty[T any]() stream.Iterator[T] { return emptyIterator[T]{} }

type emptyIterator[T any] struct{}

func (e emptyIterator[T]) TryAdvance(_ types.Consumer[T]) bool {
	return false
}

func (e emptyIterator[T]) ForEachRemaining(_ types.Consumer[T]) {
}

func (e emptyIterator[T]) Next() (T, error) {
	var zeroValue T
	return zeroValue, stream.Done
}
