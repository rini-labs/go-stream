package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/filter"
)

// Filter returns a stream consisting of the elements of this stream that match the given predicate.
func filterPipeline[OUT any](p pipeline[OUT], predicate stream.Predicate[OUT]) stream.Stream[OUT] {
	return ofPipeline(p, func(sink stream.Sink[OUT]) stream.Sink[OUT] {
		return filter.NewSink(sink, predicate)
	})
}
