package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/comparators"
	"github.com/rini-labs/go-stream/sort"
	"golang.org/x/exp/constraints"
)

func sortPipeline[OUT any](s stream.Stream[OUT], comparator stream.Comparator[OUT], naturalSort bool) stream.Stream[OUT] {
	return OfStateful(s, func(flags int, sink stream.Sink[OUT]) stream.Sink[OUT] {
		if stream.OpFlagSorted.IsKnown(flags) && naturalSort {
			return sink
		}
		return sort.NewSink(sink, comparator)
	}, stream.IsOrdered|stream.IsSorted)
}

func SortWithComparator[OUT any](s stream.Stream[OUT], comparator stream.Comparator[OUT]) stream.Stream[OUT] {
	return sortPipeline(s, comparator, false)
}

func Sort[OUT constraints.Ordered](s stream.Stream[OUT]) stream.Stream[OUT] {
	return sortPipeline(s, comparators.Natural[OUT](), true)
}
