package stream

type TerminalOp[IN any, OUT any] interface {
	evaluateSequential(helper PipelineHelper[IN], itr Iterator[IN]) OUT
}
