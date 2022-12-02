package stream

import (
	"github.com/rini-labs/go-stream/pkg/types"
)

type Sink[T any] interface {
	types.Consumer[T]

	// Begin resets the sink state to receive a fresh data set.
	Begin(size int64)

	// End indicates that all elements have been pushed.
	End()
}
