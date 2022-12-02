package iterator

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/pkg/types"
)

func New[T any](next func() (T, error)) stream.Iterator[T] {
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
