package internal

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal/iterator"
)

func Map[IT any, OT any](s stream.Stream[IT], mapper func(IT) OT) stream.Stream[OT] {
	var zeroValue OT
	iter := s.Iterator()
	return NewStreamImpl(iterator.New(func() (OT, error) {
		if nextValue, err := iter.Next(); err == nil {
			return mapper(nextValue), nil
		} else {
			return zeroValue, err
		}
	}))
}

func FlatMap[IT any, OT any](s stream.Stream[IT], mapper func(IT) stream.Stream[OT]) stream.Stream[OT] {
	var zeroValue OT

	inputStreamIter := s.Iterator()
	var nextOutputStreamIter stream.Iterator[OT]

	return NewStreamImpl(iterator.New(func() (OT, error) {
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
			case stream.Done:
				nextOutputStreamIter = nil
			default:
				return zeroValue, err
			}
		}
	}))
}
