package distinct

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal"
	"github.com/rini-labs/go-stream/internal/iterator"
)

func MakeStream[T comparable](s stream.Stream[T]) stream.Stream[T] {
	processedElems := map[T]bool{}
	iter := s.Iterator()
	return internal.NewStreamImpl(iterator.New(func() (T, error) {
		for {
			nextValue, err := iter.Next()
			if err != nil {
				return nextValue, err
			}
			if _, ok := processedElems[nextValue]; !ok {
				processedElems[nextValue] = true
				return nextValue, nil
			}
		}
	}))
}
