package stream_test

import (
	"testing"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/comparators"
	"github.com/rini-labs/go-stream/consumers"
	"github.com/rini-labs/go-stream/iterators"
	"github.com/rini-labs/go-stream/predicates"
	"github.com/rini-labs/go-stream/streams"
	"github.com/rini-labs/go-stream/supplier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	isEvent   = predicates.Of(func(value int) bool { return value%2 == 0 })
	doubleInt = func(value int) int { return value * 2 }
)

func TestStreamFromIterator(t *testing.T) {
	s := streams.Of(iterators.OfSlice([]int{1, 2, 3, 4, 5}), stream.IsSorted|stream.IsOrdered)
	assert.NotNil(t, s)

	iterator, err := s.Iterator()
	require.NoError(t, err)
	count := 0
	iterator.ForEachRemaining(consumers.Of(func(value int) {
		count++
	}))
	assert.Equal(t, 5, count)

	err = s.Close()
	require.NoError(t, err)

}

func TestStreamFromIteratorSupplier(t *testing.T) {
	s := streams.OfSupplier(supplier.Of(func() (stream.Iterator[int], error) { return iterators.OfSlice([]int{1, 2, 3, 4, 5}), nil }), stream.IsSorted|stream.IsOrdered)
	assert.NotNil(t, s)

	iterator, err := s.Iterator()
	require.NoError(t, err)
	count := 0
	iterator.ForEachRemaining(consumers.Of(func(value int) {
		count++
	}))
	assert.Equal(t, 5, count)

	err = s.Close()
	require.NoError(t, err)

}

func TestForEach(t *testing.T) {
	s := streams.Of(iterators.OfSlice([]int{1, 2, 3, 4, 5}), stream.IsSorted|stream.IsOrdered)
	assert.NotNil(t, s)

	count := 0
	err := s.ForEach(consumers.Of(func(value int) {
		count++
	}))
	require.NoError(t, err)
	assert.Equal(t, 5, count)

	err = s.Close()
	require.NoError(t, err)
}

func TestFilter(t *testing.T) {
	data, err := streams.Of(iterators.OfSlice([]int{1, 2, 3, 4, 5}), stream.IsSorted|stream.IsOrdered).Filter(isEvent).ToSlice()
	require.NoError(t, err)
	assert.Equal(t, []int{2, 4}, data)
	require.NoError(t, err)
}

func TestMapper(t *testing.T) {
	data, err := streams.Map(streams.Of(iterators.OfSlice([]int{1, 2, 3, 4, 5}), stream.IsSorted|stream.IsOrdered), doubleInt).ToSlice()
	require.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, data)
}

func TestFlatMapper(t *testing.T) {
	data, err := streams.FlatMap(streams.Of(iterators.OfSlice([]int{1, 2, 3, 4, 5}), stream.IsSorted|stream.IsOrdered), func(value int) stream.Stream[int] {
		return streams.Of(iterators.OfSlice([]int{value, value}), stream.IsSorted|stream.IsOrdered)
	}).ToSlice()
	require.NoError(t, err)
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}, data)
}

func TestDistinct(t *testing.T) {
	data, err := streams.Distinct(streams.Of(iterators.OfSlice([]int{1, 2, 1, 4, 2}), stream.IsSorted|stream.IsOrdered)).ToSlice()
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 4}, data)
}

func TestSorting(t *testing.T) {
	data, err := streams.Of(iterators.OfSlice([]int{1, 2, 1, 4, 2}), stream.IsSorted|stream.IsOrdered).Sort(comparators.Natural[int]()).ToSlice()
	require.NoError(t, err)
	assert.Equal(t, []int{1, 1, 2, 2, 4}, data)
}

func TestCount(t *testing.T) {
	count, err := streams.FlatMap(streams.Of(iterators.OfSlice([]int{1, 2, 3, 4, 5}), stream.IsSorted|stream.IsOrdered), func(value int) stream.Stream[int] {
		return streams.Of(iterators.OfSlice([]int{value, value}), stream.IsSorted|stream.IsOrdered)
	}).Count()
	require.NoError(t, err)
	assert.Equal(t, 10, count)
}
