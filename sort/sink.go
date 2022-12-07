package sort

import (
	"fmt"
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

	comparator   stream.Comparator[IN]
	data         []IN
	expectedSize int
}

func (s *sink[IN]) Begin(size int) {
	s.expectedSize = size
	if size == -1 {
		size = 8
	}
	s.data = make([]IN, 0, size)
}

func (s *sink[IN]) End() {
	if s.expectedSize != -1 && len(s.data) != s.expectedSize {
		panic(fmt.Sprintf("sort: slice length %d does not match size %d", len(s.data), s.expectedSize))
	}
	sort.Slice(s.data, func(i, j int) bool {
		return s.comparator.Compare(s.data[i], s.data[j]) < 0
	})
	sinks.CopyInto[IN](iterators.OfSlice(s.data), s.Downstream)
	s.data = nil
}

func (s *sink[IN]) Accept(t IN) {
	s.data = append(s.data, t)
}
