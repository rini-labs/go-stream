package iterator

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/pkg/types"
)

func OfSlice[T any](slice []T) stream.Iterator[T] {
	return &sliceIterator[T]{slice: slice}
}

type sliceIterator[T any] struct {
	slice []T
}

func (si *sliceIterator[T]) TryAdvance(consumer types.Consumer[T]) bool {
	if len(si.slice) == 0 {
		return false
	}

	val := si.slice[0]
	si.slice = si.slice[1:]
	consumer.Accept(val)
	return true
}

func (si *sliceIterator[T]) Next() (T, error) {
	if len(si.slice) <= 0 {
		var zeroValue T
		return zeroValue, stream.Done
	}
	rv := si.slice[0]
	si.slice = si.slice[1:]
	return rv, nil
}

func (si *sliceIterator[T]) ForEachRemaining(consumer types.Consumer[T]) {
	for si.TryAdvance(consumer) {
	}
}
