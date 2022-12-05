package sinks

import "github.com/rini-labs/go-stream"

type ChainedSink[IN any, OUT any] struct {
	Downstream stream.Sink[OUT]
}

func (cs *ChainedSink[IN, OUT]) Begin(size int) {
	cs.Downstream.Begin(size)
}

func (cs *ChainedSink[IN, OUT]) End() {
	cs.Downstream.End()
}
