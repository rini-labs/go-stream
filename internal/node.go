package internal

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal/iterator"
	"github.com/rini-labs/go-stream/pkg/types"
)

type Node[T any] interface {
	// Iterator returns an iterator for the elements of this stream.
	Iterator() stream.Iterator[T]

	ForEach(consumer types.Consumer[T])

	AsSlice() []T

	CopyInto(dst []T)

	Count() int64
}

func EmptyNode[T any]() Node[T] {
	return emptyNode[T]{}
}

type emptyNode[T any] struct{}

func (e emptyNode[T]) Iterator() stream.Iterator[T] {
	return iterator.Empty[T]()
}

func (e emptyNode[T]) ForEach(_ types.Consumer[T]) {
}

func (e emptyNode[T]) AsSlice() []T {
	return make([]T, 0)
}

func (e emptyNode[T]) CopyInto(_ []T) {
}

func (e emptyNode[T]) Count() int64 {
	return 0
}

func SliceNode[T any](slice []T) Node[T] {
	return &sliceNode[T]{slice: slice}
}

func SliceNodeOfSize[T any](size int64) Node[T] {
	return &sliceNode[T]{slice: make([]T, size), curSize: 0}
}

type sliceNode[T any] struct {
	slice   []T
	curSize int
}

func (sn *sliceNode[T]) Iterator() stream.Iterator[T] {
	if sn.curSize == 0 {
		return iterator.Empty[T]()
	}
	return iterator.OfSlice[T](sn.slice[:sn.curSize])
}

func (sn *sliceNode[T]) ForEach(consumer types.Consumer[T]) {
	for _, val := range sn.slice {
		consumer.Accept(val)
	}
}

func (sn *sliceNode[T]) AsSlice() []T {
	return sn.slice
}

func (sn *sliceNode[T]) CopyInto(dst []T) {
	if len(dst) == sn.curSize {
		panic("destination slice is not large enough")
	}
	copy(dst, sn.slice[:sn.curSize])
}

func (sn *sliceNode[T]) Count() int64 {
	return int64(sn.curSize)
}
