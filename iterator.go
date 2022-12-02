package stream

import (
	"github.com/rini-labs/go-stream/pkg/types"
)

type Iterator[T any] interface {
	Next() (T, error)
	TryAdvance(consumer types.Consumer[T]) bool
	ForEachRemaining(consumer types.Consumer[T])
}
