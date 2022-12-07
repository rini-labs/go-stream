package filter

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
)

func NewSink[OUT any](downstream stream.Sink[OUT], predicate stream.Predicate[OUT]) stream.Sink[OUT] {
	return &sink[OUT]{ChainedSink: sinks.ChainedSink[OUT, OUT]{Downstream: downstream}, predicate: predicate}
}

type sink[OUT any] struct {
	predicate stream.Predicate[OUT]
	sinks.ChainedSink[OUT, OUT]
}

func (s *sink[OUT]) Begin(_ int) {
	s.Downstream.Begin(-1)
}

func (s *sink[OUT]) Accept(val OUT) {
	if s.predicate.Test(val) {
		s.Downstream.Accept(val)
	}
}
