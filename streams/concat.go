package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/iterators"
)

func Concat[OUT any](stream1, stream2 stream.Stream[OUT]) (stream.Stream[OUT], error) {
	iter1, err := stream1.Iterator()
	if err != nil {
		return nil, err
	}
	iter2, err := stream2.Iterator()
	if err != nil {
		return nil, err
	}
	return Of(iterators.Concat(iter1, iter2), 0), nil
}
