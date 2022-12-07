package stream

type TerminalSink[IN any, OUT any] interface {
	Sink[IN]
	Supplier[OUT]
}
