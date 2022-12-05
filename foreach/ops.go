package foreach

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/sinks"
)

func Of[IN any](consumer stream.Consumer[IN]) stream.TerminalOp[IN, any] {
	return &forEachOp[IN]{consumer: consumer}
}

type forEachOp[IN any] struct {
	consumer stream.Consumer[IN]
}

func (fe *forEachOp[IN]) Accept(val IN) {
	fe.consumer.Accept(val)
}

func (fe *forEachOp[IN]) Begin(_ int) {
}

func (fe *forEachOp[IN]) End() {
}

func (fe *forEachOp[IN]) EvaluateSequential(itr stream.Iterator[IN]) any {
	sinks.CopyInto[IN](itr, fe)
	return nil
}
