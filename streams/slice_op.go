package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/slice"
)

func slicePipeline[OUT any](p pipeline[OUT], skip int, limit int) stream.Stream[OUT] {
	opsFlag := stream.NotSized
	if limit != -1 {
		opsFlag |= stream.IsShortCircuit
	}
	return OfStateful[OUT, OUT](p.(stream.Stream[OUT]), func(flags int, sink stream.Sink[OUT]) stream.Sink[OUT] {
		return slice.NewSink(sink, skip, limit)
	}, opsFlag)
}
