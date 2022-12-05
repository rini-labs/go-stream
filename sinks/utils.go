package sinks

import "github.com/rini-labs/go-stream"

func CopyInto[OUT any](iterator stream.Iterator[OUT], sink stream.Sink[OUT]) {
	sink.Begin(-1)
	iterator.ForEachRemaining(sink)
	sink.End()
}
