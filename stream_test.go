package stream_test

import (
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rini-labs/go-stream/streams"
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
	in := streams.Of(1, -1, 4, 3, 2).Sorted(func(i, j int) int { return i - j })
	slice := in.ToSlice()
	assert.Equal(t, []int{-1, 1, 2, 3, 4}, slice)
}

func TestSortedRandom(t *testing.T) {
	var source []int64
	gofakeit.Slice(&source)
	in := streams.OfSlice(source).Sorted(func(i, j int64) int {
		val := i - j
		if val < 0 {
			return -1
		} else if val > 0 {
			return 1
		}
		return 0
	})
	slice := in.ToSlice()
	assert.True(t, sort.SliceIsSorted(slice, func(i, j int) bool { return slice[i] < slice[j] }))
}
