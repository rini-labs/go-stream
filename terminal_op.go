package stream

type TerminalOp[IN any, OUT any] interface {
	EvaluateSequential(helper PipelineHelper[IN], itr Iterator[IN]) OUT
}
