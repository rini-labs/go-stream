package peek

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
)

func NewSink[IN any](downstream stream.Sink[IN], consumer stream.Consumer[IN]) stream.Sink[IN] {
	return &sink[IN]{ChainedSink: sinks.ChainedSink[IN, IN]{Downstream: downstream}, consumer: consumer}
}

type sink[IN any] struct {
	sinks.ChainedSink[IN, IN]

	consumer stream.Consumer[IN]
}

func (s *sink[IN]) Accept(t IN) {
	s.consumer.Accept(t)
	s.Downstream.Accept(t)
}
