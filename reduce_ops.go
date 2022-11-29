package stream

type AccumulatingSink[IN any, OUT any] interface {
	TerminalSink[IN, OUT]
	Combine(other AccumulatingSink[IN, OUT])
}

type box[T any] struct {
	state T
}

func (b *box[T]) Get() T {
	return b.state
}

type CountingSink[T any] struct {
	abstractSink[T]
	box[int64]
}

func (cs *CountingSink[T]) Get() int64 {
	return cs.state
}

func (cs *CountingSink[T]) Combine(other AccumulatingSink[T, int64]) {
	cs.state += other.Get()
}

func (cs *CountingSink[T]) Accept(_ T) {
	cs.state++
}

func (cs *CountingSink[T]) Begin(_ int64) {
	cs.state = 0
}

var _ AccumulatingSink[int, int64] = (*CountingSink[int])(nil)

type SinkProvider[IN any, OUT any] func() AccumulatingSink[IN, OUT]

type ReduceOps[IN any, OUT any] interface {
	TerminalOp[IN, OUT]
}

func NewReduceOp[IN any, OUT any](sinkProvider SinkProvider[IN, OUT]) ReduceOps[IN, OUT] {
	return &reduceOp[IN, OUT]{sinkProvider: sinkProvider}
}

type reduceOp[IN any, OUT any] struct {
	sinkProvider SinkProvider[IN, OUT]
}

func (ro *reduceOp[IN, OUT]) evaluateSequential(helper PipelineHelper[IN], itr Iterator[IN]) OUT {
	return helper.WrapAndCopyInto(ro.sinkProvider(), itr).(AccumulatingSink[IN, OUT]).Get()
}

var _ ReduceOps[int, int64] = (*reduceOp[int, int64])(nil)

func MakeCounting[T any]() TerminalOp[T, int64] {
	return NewReduceOp[T, int64](func() AccumulatingSink[T, int64] { return &CountingSink[T]{} })
}
