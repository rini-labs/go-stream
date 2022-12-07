package nodes

import (
	"errors"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/iterators"
)

func NewSliceBuilder[T any](size int) NodeBuilder[T] {
	return &sliceNodeBuilder[T]{sliceNode: sliceNode[T]{data: make([]T, 0, size), size: size}}
}

type sliceNode[T any] struct {
	data []T
	size int
}

// Iterator implements Node
func (sn *sliceNode[T]) Iterator() (stream.Iterator[T], error) {
	if sn.size != len(sn.data) {
		return nil, errors.New("illegal state: slice size does not match data length")
	}
	return iterators.OfSlice(sn.data), nil
}

// ForEach implements Node
func (sn *sliceNode[T]) ForEach(consumer stream.Consumer[T]) error {
	if sn.size != len(sn.data) {
		return errors.New("illegal state: slice size does not match data length")
	}
	for _, value := range sn.data {
		consumer.Accept(value)
	}

	return nil
}

func (sn *sliceNode[T]) AsArray() ([]T, error) {
	if sn.size != len(sn.data) {
		return nil, errors.New("illegal state: slice size does not match data length")
	}
	return sn.data, nil
}

// Count implements Node
func (sn *sliceNode[T]) Count() int {
	return len(sn.data)
}

type sliceNodeBuilder[T any] struct {
	sliceNode[T]
}

func (snb *sliceNodeBuilder[T]) Build() Node[T] {
	return snb
}

func (snb *sliceNodeBuilder[T]) Begin(size int) {
	if size != snb.size {
		panic("illegal state: slice size does not match data length")
	}
	snb.data = make([]T, 0, size)
}

func (snb *sliceNodeBuilder[T]) Accept(t T) {
	snb.data = append(snb.data, t)
}

func (snb *sliceNodeBuilder[T]) End() {
	if len(snb.data) != snb.size {
		panic("illegal state: slice size does not match data length")
	}
}
