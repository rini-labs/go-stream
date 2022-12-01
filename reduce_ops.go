package stream

import "github.com/rini-labs/go-stream/types"

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

type reduceSink[T any, U any] struct {
	box[U]
	reducer  types.BiFunction[U, T, U]
	combiner types.BinaryOperator[U]
	seed     U
}

func (rs *reduceSink[T, U]) Begin(_ int64) {
	rs.state = rs.seed
}

func (rs *reduceSink[T, U]) Accept(val T) {
	rs.state = rs.reducer.Apply(rs.state, val)
}

func (rs *reduceSink[T, U]) End() {
}

func (rs *reduceSink[T, U]) Combine(other AccumulatingSink[T, U]) {
	otherSink := other.(*reduceSink[T, U])
	rs.state = rs.combiner.Apply(rs.state, otherSink.state)
}

func MakeReducer[T any, U any](seed U, reducer types.BiFunction[U, T, U], combiner types.BinaryOperator[U]) TerminalOp[T, U] {
	return NewReduceOp[T, U](func() AccumulatingSink[T, U] {
		return &reduceSink[T, U]{reducer: reducer, combiner: combiner, seed: seed}
	})
}

func MakeReducerFromOperator[T any](operator types.BinaryOperator[T]) TerminalOp[T, types.Optional[T]] {
	reducer := func(state types.Optional[T], val T) types.Optional[T] {
		if !state.IsPresent() {
			return types.OptionalOf[T](val)
		}
		stateVal, _ := state.Get()
		return types.OptionalOf[T](operator.Apply(stateVal, val))
	}
	combiner := func(state types.Optional[T], otherState types.Optional[T]) types.Optional[T] {
		if !otherState.IsPresent() {
			return state
		}
		otherStateVal, _ := otherState.Get()
		return reducer(state, otherStateVal)
	}

	return MakeReducer[T, types.Optional[T]](types.EmptyOptional[T](), types.BiFunctionFromFunc(reducer), types.BinaryOperatorFromFunc(combiner))
}
