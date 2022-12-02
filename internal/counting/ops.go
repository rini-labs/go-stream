package counting

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal/reduce"
	"github.com/rini-labs/go-stream/pkg/types"
)

type countingSink[T any] struct {
	state int64
}

func (cs *countingSink[T]) Get() int64 {
	return cs.state
}

func (cs *countingSink[T]) Combine(other reduce.AccumulatingSink[T, int64]) {
	cs.state += other.Get()
}

func (cs *countingSink[T]) Accept(_ T) {
	cs.state++
}

func (cs *countingSink[T]) Begin(_ int64) {
	cs.state = 0
}

func (cs *countingSink[T]) End() {
}

var _ reduce.AccumulatingSink[int, int64] = (*countingSink[int])(nil)

func MakeCounting[T any]() stream.TerminalOp[T, int64] {
	sinkSupplier := types.NewSupplier(func() reduce.AccumulatingSink[T, int64] { return &countingSink[T]{} })
	return reduce.NewReduceOp[T, int64](sinkSupplier)
}
