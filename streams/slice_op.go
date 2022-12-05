package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/slice"
)

func slicePipeline[OUT any](p pipeline[OUT], skip int, limit int) *derivedPipeline[OUT, OUT] {
	return ofPipeline(p, func(sink stream.Sink[OUT]) stream.Sink[OUT] {
		return slice.NewSink(sink, skip, limit)
	})
}
