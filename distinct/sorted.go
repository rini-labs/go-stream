package distinct

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
)

func NewSinkSorted[IN comparable](downstream stream.Sink[IN]) stream.Sink[IN] {
	return &sinkSorted[IN]{ChainedSink: sinks.ChainedSink[IN, IN]{Downstream: downstream}}
}

type sinkSorted[IN comparable] struct {
	sinks.ChainedSink[IN, IN]

	seenNull bool
	lastSeen IN
}

func (s *sinkSorted[IN]) Begin(size int) {
	var zeroVal IN
	s.seenNull = false
	s.lastSeen = zeroVal
	s.Downstream.Begin(-1)
}

func (s *sinkSorted[IN]) End() {
	var zeroVal IN
	s.seenNull = false
	s.lastSeen = zeroVal
	s.Downstream.End()
}

func (s *sinkSorted[IN]) Accept(t IN) {
	var zeroVal IN
	if t == zeroVal {
		if !s.seenNull {
			s.seenNull = true
			s.lastSeen = zeroVal
			s.Downstream.Accept(t)
		}
	} else if s.lastSeen == zeroVal || t != s.lastSeen {
		s.lastSeen = t
		s.Downstream.Accept(t)
	}
}
