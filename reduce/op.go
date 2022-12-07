package reduce

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
	"github.com/rini-labs/go-stream/supplier"
)

func NewOpWithSeed[IN any, OUT any](seed OUT, reducer func(OUT, IN) OUT) stream.TerminalOp[IN, OUT] {
	return &reduceOp[IN, OUT, *reducingSink[IN, OUT]]{sinkSupplier: supplier.Of(func() (*reducingSink[IN, OUT], error) {
		return &reducingSink[IN, OUT]{seed: seed, reducer: reducer}, nil
	})}
}

func NewOp[IN any, OUT any](reducer func(OUT, IN) OUT) stream.TerminalOp[IN, OUT] {
	return &reduceOp[IN, OUT, *reducingSink[IN, OUT]]{sinkSupplier: supplier.Of(func() (*reducingSink[IN, OUT], error) {
		return &reducingSink[IN, OUT]{reducer: reducer}, nil
	})}
}

type reduceOp[IN any, OUT any, S AccumulatingSink[IN, OUT]] struct {
	sinkSupplier stream.Supplier[S]
}

func (rop *reduceOp[IN, OUT, S]) EvaluateSequential(itr stream.Iterator[IN]) OUT {
	s, err := rop.sinkSupplier.Get()
	if err != nil {
		panic(err)
	}
	sinks.CopyInto[IN](itr, s)
	rv, err := s.Get()
	if err != nil {
		panic(err)
	}
	return rv
}
