package for_each

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/pkg/types"
)

func MakeForEach[IN any](consumer types.Consumer[IN]) stream.TerminalOp[IN, any] {
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

func (fe *forEachOp[IN]) EvaluateSequential(helper stream.PipelineHelper[IN], itr stream.Iterator[IN]) any {
	return helper.WrapAndCopyInto(fe, itr).(stream.TerminalSink[IN, any]).Get()
}
