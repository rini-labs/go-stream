package distinct

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
)

func NewSink[IN comparable](downstream stream.Sink[IN]) stream.Sink[IN] {
	return &sink[IN]{ChainedSink: sinks.ChainedSink[IN, IN]{Downstream: downstream}}
}

type sink[IN comparable] struct {
	sinks.ChainedSink[IN, IN]

	alreadyProcessed map[IN]struct{}
}

func (s *sink[IN]) Begin(size int) {
	if size == -1 {
		size = 0
	}
	s.alreadyProcessed = make(map[IN]struct{}, size)
	s.Downstream.Begin(-1)
}

func (s *sink[IN]) End() {
	s.alreadyProcessed = nil
	s.Downstream.End()
}

func (s *sink[IN]) Accept(t IN) {
	if _, ok := s.alreadyProcessed[t]; !ok {
		s.alreadyProcessed[t] = struct{}{}
		s.Downstream.Accept(t)
	}
}
