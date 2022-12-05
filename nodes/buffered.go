package nodes

import (
	"errors"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/iterators"
)

func Builder[T any]() NodeBuilder[T] {
	return &bufferedNodeBuilder[T]{Buffer: iterators.NewBuffer[T]()}
}

type bufferedNodeBuilder[T any] struct {
	iterators.Buffer[T]
	building bool
}

// Begin implements Buffer
func (bn *bufferedNodeBuilder[OUT]) Begin(size int) {
	if bn.building {
		panic("illegal state: already building")
	}
	bn.building = true
	bn.Buffer.Clear()
	bn.Buffer.Begin(size)
}

// End implements Buffer
func (bn *bufferedNodeBuilder[OUT]) End() {
	if !bn.building {
		panic("illegal state: End() called without Begin()")
	}
	bn.Buffer.End()
	bn.building = false
}

// Accept implements Buffer
func (bn *bufferedNodeBuilder[OUT]) Accept(val OUT) {
	if !bn.building {
		panic("illegal state: End() called without Begin()")
	}
	bn.Buffer.Accept(val)
}

// Iterator implements Buffer
func (bn *bufferedNodeBuilder[OUT]) Iterator() (stream.Iterator[OUT], error) {
	if bn.building {
		panic("illegal state: during building")
	}
	return bn.Buffer.Iterator()
}

func (bn *bufferedNodeBuilder[OUT]) ForEach(consumer stream.Consumer[OUT]) error {
	if bn.building {
		panic("illegal state: during building")
	}
	return bn.Buffer.ForEach(consumer)
}

// Count implements Buffer
func (bn *bufferedNodeBuilder[OUT]) Count() int {
	if bn.building {
		panic("illegal state: during building")
	}
	return bn.Buffer.Count()
}

func (bn *bufferedNodeBuilder[OUT]) Get(index int) (OUT, error) {
	if bn.building {
		panic("illegal state: during building")
	}
	return bn.Buffer.Get(index)
}

func (bn *bufferedNodeBuilder[OUT]) Build() Node[OUT] {
	if bn.building {
		panic("illegal state: during building")
	}
	return bn
}

func (bn *bufferedNodeBuilder[T]) AsArray() ([]T, error) {
	if bn.building {
		return nil, errors.New("illegal state: during building")
	}
	size := bn.Count()
	data := make([]T, size)
	bn.CopyInto(data)
	return data, nil
}
