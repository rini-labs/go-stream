package stream

import (
	"github.com/rini-labs/go-stream/pkg/types"
)

// Stream is a sequence of elements
type Stream[T any] interface {
	// Close closes the stream
	Close() error

	// Iterator returns an iterator for the elements of this stream.
	Iterator() Iterator[T]

	// OnClose returns an equivalent stream with an additional close handler.
	// Close handlers are eun when the close() method is called on the stream,
	// and are executed in the order they were added.
	// All close handlers are run even if earlier close handlers return error.
	// If any close handlers return error the first error wil be relayed to the caller.
	OnClose(closeHandler func() error) Stream[T]

	// Filter returns a stream consisting of the elements of this stream that
	// match the given predicate
	Filter(predicate types.Predicate[T]) Stream[T]

	// Map returns a stream consisting of the results of applying the given
	// function to the elements of this stream.
	// For changing the stream element type, use the stream.Map method.
	Map(mapper func(T) T) Stream[T]

	// FlatMap returns a stream consisting of the results of replacing each
	// element of this stream with the contents of a mapped stream produced
	// by applying the provided mapping function to each element.
	// Each mapped stream is closed after its contents have been placed into
	// this stream. (If a mapped stream is null an empty stream is used, instead.)
	FlatMap(func(T) Stream[T]) Stream[T]

	Peek(consumer func(T)) Stream[T]

	Limit(int64) Stream[T]

	Sorted(comparator types.Comparator[T]) Stream[T]

	Skip(count int64) Stream[T]

	ForEach(consumer types.Consumer[T])
	Count() int64

	Reduce(accumulator types.BinaryOperator[T]) types.Optional[T]

	// ToSlice returns a slice containing the elements of this stream.
	ToSlice() []T
}

func Filter[T any](input Stream[T], predicate func(T) bool) Stream[T] {
	return input.Filter(predicate)
}

func Peek[T any](s Stream[T], consumer func(T)) Stream[T] {
	return s.Peek(consumer)
}

func Limit[T any](s Stream[T], maxSize int64) Stream[T] {
	return s.Limit(maxSize)
}

func Sorted[T any](s Stream[T], comparator types.Comparator[T]) Stream[T] {
	return s.Sorted(comparator)
}

func Skip[T any](s Stream[T], count int64) Stream[T] {
	return s.Skip(count)
}

func Count[T any](s Stream[T]) int64 {
	return s.Count()
}
