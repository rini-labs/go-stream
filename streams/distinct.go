package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/distinct"
)

func distinctPipeline[IN comparable](p pipeline[IN]) stream.Stream[IN] {
	return ofPipeline(p, func(sink stream.Sink[IN]) stream.Sink[IN] {
		return distinct.NewSink(sink)
	})
}

func Distinct[IN comparable](s stream.Stream[IN]) stream.Stream[IN] {
	return distinctPipeline(s.(pipeline[IN]))
}
