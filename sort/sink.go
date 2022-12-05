package sort

import (
	"sort"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/iterators"
	"github.com/rini-labs/go-stream/sinks"
)

func NewSink[IN any](downstream stream.Sink[IN], comparator stream.Comparator[IN]) stream.Sink[IN] {
	return &sink[IN]{ChainedSink: sinks.ChainedSink[IN, IN]{Downstream: downstream}, comparator: comparator}
}

type sink[IN any] struct {
	sinks.ChainedSink[IN, IN]

	comparator stream.Comparator[IN]
	data       []IN
}

func (s *sink[IN]) Begin(size int) {
	if size == -1 {
		size = 0
	}
	s.data = make([]IN, size)
}

func (s *sink[IN]) End() {
	sort.Slice(s.data, func(i, j int) bool {
		return s.comparator.Compare(s.data[i], s.data[j]) < 0
	})
	sinks.CopyInto[IN](iterators.OfSlice(s.data), s.Downstream)
	s.data = nil
}

func (s *sink[IN]) Accept(t IN) {
	s.data = append(s.data, t)
}
