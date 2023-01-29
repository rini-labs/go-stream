package streams_test

import (
	"testing"

	"github.com/rini-labs/go-stream/iterators"
	"github.com/rini-labs/go-stream/streams"
	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {
	s1 := streams.Of(iterators.OfSlice([]int{1, 2, 3}), 0)
	s2 := streams.Of(iterators.OfSlice([]int{4, 5, 6}), 0)
	concatStream, err := streams.Concat(s1, s2)
	require.NoError(t, err)
	output, err := concatStream.ToSlice()
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6}, output)
}
