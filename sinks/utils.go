package sinks

import "github.com/rini-labs/go-stream"

func CopyInto[OUT any](iterator stream.Iterator[OUT], sink stream.Sink[OUT]) {
	sink.Begin(iterator.GetExactSizeIfKnown())
	iterator.ForEachRemaining(sink)
	sink.End()
}
