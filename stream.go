package stream

type Predicate[T any] func(T) bool

type Iterator[T any] interface {
	Next() (T, error)
}

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
	Filter(predicate Predicate[T]) Stream[T]

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

	// ToSlice returns a slice containing the elements of this stream.
	ToSlice() []T
}

func Map[IT any, OT any](s Stream[IT], mapper func(IT) OT) Stream[OT] {
	var zeroValue OT
	iter := s.Iterator()
	return NewStreamImpl(NewIterator(func() (OT, error) {
		if nextValue, err := iter.Next(); err == nil {
			return mapper(nextValue), nil
		} else {
			return zeroValue, err
		}
	}))
}

func FlatMap[IT any, OT any](s Stream[IT], mapper func(IT) Stream[OT]) Stream[OT] {
	var zeroValue OT

	inputStreamIter := s.Iterator()
	var nextOutputStreamIter Iterator[OT]

	return NewStreamImpl(NewIterator(func() (OT, error) {
		for {
			if nextOutputStreamIter == nil {
				if nextValue, err := inputStreamIter.Next(); err != nil {
					nextOutputStreamIter = mapper(nextValue).Iterator()
				} else {
					return zeroValue, err
				}
			}
			nextValue, err := nextOutputStreamIter.Next()
			switch err {
			case nil:
				return nextValue, nil
			case Done:
				nextOutputStreamIter = nil
			default:
				return zeroValue, err
			}
		}
	}))
}
