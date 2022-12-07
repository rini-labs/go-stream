package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/distinct"
)

func distinctPipeline[IN comparable](p pipeline[IN]) stream.Stream[IN] {
	return OfStateful(p, func(flags int, sink stream.Sink[IN]) stream.Sink[IN] {
		if stream.OpFlagDistinct.IsKnown(flags) {
			return sink
		} else if stream.OpFlagSorted.IsKnown(flags) {
			return distinct.NewSinkSorted(sink)
		}
		return distinct.NewSink(sink)
	}, stream.IsDistinct|stream.NotSized)
}

func Distinct[IN comparable](s stream.Stream[IN]) stream.Stream[IN] {
	return distinctPipeline(s.(pipeline[IN]))
}
