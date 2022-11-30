package stream

import "github.com/rini-labs/go-stream/types"

type Iterator[T any] interface {
	Next() (T, error)
	TryAdvance(consumer types.Consumer[T]) bool
	ForEachRemaining(consumer types.Consumer[T])
}

func EmptyIterator[T any]() Iterator[T] { return emptyIterator[T]{} }

type emptyIterator[T any] struct{}

func (e emptyIterator[T]) TryAdvance(consumer types.Consumer[T]) bool {
	return false
}

func (e emptyIterator[T]) ForEachRemaining(consumer types.Consumer[T]) {
}

func (e emptyIterator[T]) Next() (T, error) {
	var zeroValue T
	return zeroValue, Done
}

func NewIterator[T any](next func() (T, error)) Iterator[T] {
	return &iterator[T]{next: next}
}

type iterator[T any] struct {
	next func() (T, error)
}

func (si *iterator[T]) Next() (T, error) {
	return si.next()
}

func (si *iterator[T]) ForEachRemaining(consumer types.Consumer[T]) {
	for si.TryAdvance(consumer) {
	}
}

func (si *iterator[T]) TryAdvance(action types.Consumer[T]) bool {
	next, err := si.next()
	if err != nil {
		return false
	}
	action.Accept(next)
	return true
}

func SliceIterator[T any](slice []T) Iterator[T] { return &sliceIterator[T]{slice: slice} }

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
		return zeroValue, Done
	}
	rv := si.slice[0]
	si.slice = si.slice[1:]
	return rv, nil
}

func (si *sliceIterator[T]) ForEachRemaining(consumer types.Consumer[T]) {
	for si.TryAdvance(consumer) {
	}
}
