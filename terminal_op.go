package stream

type TerminalOp[IN any, OUT any] interface {
	EvaluateSequential(itr Iterator[IN]) OUT
}
