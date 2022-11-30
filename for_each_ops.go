package stream

import "github.com/rini-labs/go-stream/types"

func MakeForEach[IN any](consumer types.Consumer[IN]) TerminalOp[IN, any] {
	return &forEachOp[IN]{consumer: consumer}
}

type forEachOp[IN any] struct {
	consumer types.Consumer[IN]
}

func (fe *forEachOp[IN]) Accept(val IN) {
	fe.consumer.Accept(val)
}

func (fe *forEachOp[IN]) Begin(_ int64) {
}

func (fe *forEachOp[IN]) End() {
}

func (fe *forEachOp[IN]) Get() any {
	return nil
}

func (fe *forEachOp[IN]) evaluateSequential(helper PipelineHelper[IN], itr Iterator[IN]) any {
	return helper.WrapAndCopyInto(fe, itr).(TerminalSink[IN, any]).Get()
}
