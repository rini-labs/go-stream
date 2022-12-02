package transform

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal"
)

func Map[IT any, OT any](s stream.Stream[IT], mapper func(IT) OT) stream.Stream[OT] {
	return internal.Map[IT, OT](s, mapper)
}

func FlatMap[IT any, OT any](s stream.Stream[IT], mapper func(IT) stream.Stream[OT]) stream.Stream[OT] {
	return internal.FlatMap[IT, OT](s, mapper)
}
