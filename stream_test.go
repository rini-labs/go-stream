package stream_test

import (
	"strconv"
	"testing"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/streams"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	in := streams.Of[int](1, 2, 3, 4, 5)
	double := in.Map(func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, double.ToSlice())

	strings := stream.Map(in, strconv.Itoa)
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, strings.ToSlice())
}
