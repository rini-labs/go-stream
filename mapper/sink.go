package mapper

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
)

func NewMapSink[IN any, OUT any](downstream stream.Sink[OUT], mapper func(v IN) OUT) stream.Sink[IN] {
	return &mapperSink[IN, OUT]{ChainedSink: sinks.ChainedSink[IN, OUT]{Downstream: downstream}, mapper: mapper}
}

type mapperSink[IN any, OUT any] struct {
	mapper func(v IN) OUT
	sinks.ChainedSink[IN, OUT]
}

func (s *mapperSink[IN, OUT]) Accept(val IN) {
	s.Downstream.Accept(s.mapper(val))
}

func NewFlatMapSink[IN any, OUT any](downstream stream.Sink[OUT], mapper func(v IN) stream.Stream[OUT]) stream.Sink[IN] {
	return &flatMapperSink[IN, OUT]{ChainedSink: sinks.ChainedSink[IN, OUT]{Downstream: downstream}, mapper: mapper}
}

type flatMapperSink[IN any, OUT any] struct {
	mapper func(v IN) stream.Stream[OUT]
	sinks.ChainedSink[IN, OUT]
}

func (s *flatMapperSink[IN, OUT]) Begin(_ int) {
	s.Downstream.Begin(-1)
}

func (s *flatMapperSink[IN, OUT]) Accept(val IN) {
	if st := s.mapper(val); st != nil {
		st.ForEach(s.Downstream)
	}
}
