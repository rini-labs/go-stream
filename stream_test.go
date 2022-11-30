package stream_test

import (
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rini-labs/go-stream/comparators"
	"github.com/rini-labs/go-stream/streams"
	"github.com/rini-labs/go-stream/types"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	in := streams.Of[int](1, 2, 3, 4, 5)
	double := in.Map(func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, double.ToSlice())
}

func TestLimit(t *testing.T) {
	in := streams.Of[int](1, 2, 3, 4, 5).Limit(2)
	assert.Equal(t, []int{1, 2}, in.ToSlice())
}

func TestSkip(t *testing.T) {
	in := streams.Of[int](1, 2, 3, 4, 5).Skip(2)
	assert.Equal(t, []int{3, 4, 5}, in.ToSlice())
}

func TestSorted(t *testing.T) {
	out := streams.Of(1, -1, 4, 3, 2).Sorted(comparators.Int[int]).ToSlice()
	assert.Equal(t, []int{-1, 1, 2, 3, 4}, out)
}

func TestSortedRandom(t *testing.T) {
	var source []int
	gofakeit.Slice(&source)
	out := streams.OfSlice(source).Sorted(comparators.Int[int]).ToSlice()
	assert.True(t, sort.SliceIsSorted(out, func(i, j int) bool { return out[i] < out[j] }))
}

func TestCount(t *testing.T) {
	count := streams.Of[int](1, 2, 3, 4, 5).Count()
	assert.EqualValues(t, 5, count)
}

func TestForEach(t *testing.T) {
	count := 0
	streams.Of[int](1, 2, 3, 4, 5).ForEach(types.ToConsumer(func(i int) {
		count++
	}))
	assert.EqualValues(t, 5, count)
}
