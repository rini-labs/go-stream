package stream

import (
	"github.com/rini-labs/go-stream/pkg/types"
)

type TerminalSink[IN any, OUT any] interface {
	Sink[IN]
	types.Supplier[OUT]
}
