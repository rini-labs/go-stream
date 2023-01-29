package iterators_test

import (
	"testing"

	"github.com/rini-labs/go-stream/iterators"
	"github.com/rini-labs/go-stream/streams"
	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {
	iter1 := iterators.OfSlice([]int{1, 2, 3})
	iter2 := iterators.OfSlice([]int{4, 5, 6})
	concatIter := iterators.Concat(iter1, iter2)
	output, err := streams.Of(concatIter, 0).ToSlice()
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, output)
}
