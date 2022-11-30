package stream

type countingSink[T any] struct {
	box[int64]
}

func (cs *countingSink[T]) Get() int64 {
	return cs.state
}

func (cs *countingSink[T]) Combine(other AccumulatingSink[T, int64]) {
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

var _ AccumulatingSink[int, int64] = (*countingSink[int])(nil)

func MakeCounting[T any]() TerminalOp[T, int64] {
	return NewReduceOp[T, int64](func() AccumulatingSink[T, int64] { return &countingSink[T]{} })
}
