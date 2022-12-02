package reduce

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/pkg/types"
)

type Ops[IN any, OUT any] interface {
	stream.TerminalOp[IN, OUT]
}

type SinkFactory[IN any, OUT any] func() AccumulatingSink[IN, OUT]

func NewReduceOp[IN any, OUT any](sinkSupplier types.Supplier[AccumulatingSink[IN, OUT]]) Ops[IN, OUT] {
	return &reduceOp[IN, OUT]{sinkSupplier: sinkSupplier}
}

type reduceOp[IN any, OUT any] struct {
	sinkSupplier types.Supplier[AccumulatingSink[IN, OUT]]
}

func (ro *reduceOp[IN, OUT]) EvaluateSequential(helper stream.PipelineHelper[IN], itr stream.Iterator[IN]) OUT {
	return helper.WrapAndCopyInto(ro.sinkSupplier.Get(), itr).(AccumulatingSink[IN, OUT]).Get()
}

type sink[T any, U any] struct {
	state    U
	reducer  types.BiFunction[U, T, U]
	combiner types.BinaryOperator[U]
	seed     U
}

func (rs *sink[T, U]) Get() U {
	return rs.state
}

func (rs *sink[T, U]) Begin(_ int64) {
	rs.state = rs.seed
}

func (rs *sink[T, U]) Accept(val T) {
	rs.state = rs.reducer.Apply(rs.state, val)
}

func (rs *sink[T, U]) End() {
}

func (rs *sink[T, U]) Combine(other AccumulatingSink[T, U]) {
	otherSink := other.(*sink[T, U])
	rs.state = rs.combiner.Apply(rs.state, otherSink.state)
}

func MakeReducer[T any, U any](seed U, reducer types.BiFunction[U, T, U], combiner types.BinaryOperator[U]) stream.TerminalOp[T, U] {
	sinkSupplier := types.NewSupplier(func() AccumulatingSink[T, U] {
		return &sink[T, U]{reducer: reducer, combiner: combiner, seed: seed}
	})
	return NewReduceOp[T, U](sinkSupplier)
}

func MakeReducerFromOperator[T any](operator types.BinaryOperator[T]) stream.TerminalOp[T, types.Optional[T]] {
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
