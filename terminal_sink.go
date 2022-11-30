package stream

import "github.com/rini-labs/go-stream/types"

type TerminalSink[IN any, OUT any] interface {
	Sink[IN]
	types.Supplier[OUT]
}
