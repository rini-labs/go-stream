package slice

import (
	"math"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
	"golang.org/x/exp/constraints"
)

func NewSink[IN any](downstream stream.Sink[IN], skip int, limit int) stream.Sink[IN] {
	if limit < 0 {
		limit = math.MaxInt
	}
	return &sink[IN]{ChainedSink: sinks.ChainedSink[IN, IN]{Downstream: downstream}, skip: skip, limit: limit}
}

type sink[IN any] struct {
	sinks.ChainedSink[IN, IN]

	skip  int
	limit int
}

func (s *sink[IN]) Begin(size int) {
	if size >= 0 {
		size = max(-1, min(size-s.skip, s.limit))
	}
	s.Downstream.Begin(size)
}

func (s *sink[IN]) Accept(t IN) {
	if s.skip == 0 {
		if s.limit > 0 {
			s.limit--
			s.Downstream.Accept(t)
		}
	} else {
		s.skip--
	}
}

func max[T constraints.Integer](x, y T) T {
	if x < y {
		return y
	}
	return x
}

func min[T constraints.Integer](x, y T) T {
	if x > y {
		return y
	}
	return x
}
