package reduce

import "github.com/rini-labs/go-stream"

type AccumulatingSink[IN any, OUT any] interface {
	stream.TerminalSink[IN, OUT]
	Combine(other AccumulatingSink[IN, OUT])
}
