package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/peek"
)

func peekPipeline[OUT any](p pipeline[OUT], consumer stream.Consumer[OUT]) stream.Stream[OUT] {
	return ofPipeline(p, func(sink stream.Sink[OUT]) stream.Sink[OUT] {
		return peek.NewSink(sink, consumer)
	})
}
