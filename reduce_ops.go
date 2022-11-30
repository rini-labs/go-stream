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

type ReduceOps[IN any, OUT any] interface {
	TerminalOp[IN, OUT]
}

type SinkFactory[IN any, OUT any] func() AccumulatingSink[IN, OUT]

func NewReduceOp[IN any, OUT any](sinkProvider SinkFactory[IN, OUT]) ReduceOps[IN, OUT] {
	return &reduceOp[IN, OUT]{makeSink: sinkProvider}
}

type reduceOp[IN any, OUT any] struct {
	makeSink SinkFactory[IN, OUT]
}

func (ro *reduceOp[IN, OUT]) evaluateSequential(helper PipelineHelper[IN], itr Iterator[IN]) OUT {
	return helper.WrapAndCopyInto(ro.makeSink(), itr).(AccumulatingSink[IN, OUT]).Get()
}

var _ ReduceOps[int, int64] = (*reduceOp[int, int64])(nil)
