package reduce

import "github.com/rini-labs/go-stream"

type AccumulatingSink[IN any, OUT any] interface {
	stream.TerminalSink[IN, OUT]
}

type reducingSink[IN any, OUT any] struct {
	state OUT

	seed    OUT
	reducer func(OUT, IN) OUT
}

func (rs *reducingSink[IN, OUT]) Accept(value IN) {
	rs.state = rs.reducer(rs.state, value)
}

func (rs *reducingSink[IN, OUT]) Begin(_ int) {
	rs.state = rs.seed
}

func (rs *reducingSink[IN, OUT]) End() {
}

func (rs *reducingSink[IN, OUT]) Get() (OUT, error) {
	return rs.state, nil
}
