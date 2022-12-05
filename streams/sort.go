package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sort"
)

func sortPipeline[OUT any](p pipeline[OUT], comparator stream.Comparator[OUT]) *derivedPipeline[OUT, OUT] {
	return ofPipeline(p, func(sink stream.Sink[OUT]) stream.Sink[OUT] {
		return sort.NewSink(sink, comparator)
	})
}
