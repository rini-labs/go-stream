package internal

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal/iterator"
	"github.com/rini-labs/go-stream/pkg/types"
)

type NodeBuilder[T any] interface {
	stream.Sink[T]
	Build() Node[T]
}

func MakeNodeBuilder[T any](exactKnownSize int) NodeBuilder[T] {
	if exactKnownSize >= 0 {
		return FixedNodeBuilder[T](exactKnownSize)
	}
	return SpinedNodeBuilder[T]()
}

func FixedNodeBuilder[T any](size int) NodeBuilder[T] {
	return &fixedNodeBuilder[T]{sliceNode: sliceNode[T]{slice: make([]T, size), curSize: 0}}
}

type fixedNodeBuilder[T any] struct {
	sliceNode[T]
}

func (fb *fixedNodeBuilder[T]) Build() Node[T] {
	if fb.curSize != len(fb.slice) {
		panic("fixedNodeBuilder: size mismatch")
	}
	return fb
}

func (fb *fixedNodeBuilder[T]) Begin(size int64) {
	if int(size) != len(fb.slice) {
		panic("fixedNodeBuilder: size mismatch")
	}
	fb.curSize = 0
}

func (fb *fixedNodeBuilder[T]) Accept(t T) {
	if fb.curSize >= len(fb.slice) {
		panic("fixedNodeBuilder: size mismatch")
	}
	fb.slice[fb.curSize] = t
	fb.curSize++
}

func (fb *fixedNodeBuilder[T]) End() {
	if fb.curSize != len(fb.slice) {
		panic("fixedNodeBuilder: size mismatch")
	}
}

func SpinedNodeBuilder[T any]() NodeBuilder[T] {
	return &spinedNodeBuilder[T]{slice: make([]T, 0)}
}

type spinedNodeBuilder[T any] struct {
	slice    []T
	building bool
}

func (sb *spinedNodeBuilder[T]) Build() Node[T] {
	if sb.building {
		panic("spinedNodeBuilder: building")
	}
	return sb
}

func (sb *spinedNodeBuilder[T]) Begin(size int64) {
	if sb.building {
		panic("spinedNodeBuilder: already building")
	}
	sb.building = true
	if cap(sb.slice) < int(size) {
		oldSlice := sb.slice
		sb.slice = make([]T, len(oldSlice), size)
		copy(sb.slice, oldSlice)
	}
}

func (sb *spinedNodeBuilder[T]) Accept(t T) {
	if !sb.building {
		panic("spinedNodeBuilder: not building")
	}
	sb.slice = append(sb.slice, t)
}

func (sb *spinedNodeBuilder[T]) End() {
	if !sb.building {
		panic("spinedNodeBuilder: was not building")
	}
	sb.building = false
}

func (sb *spinedNodeBuilder[T]) Iterator() stream.Iterator[T] {
	if sb.building {
		panic("spinedNodeBuilder: still building")
	}
	if len(sb.slice) == 0 {
		return iterator.Empty[T]()
	}
	return iterator.OfSlice[T](sb.slice)
}

func (sb *spinedNodeBuilder[T]) ForEach(consumer types.Consumer[T]) {
	if sb.building {
		panic("spinedNodeBuilder: still building")
	}
	for _, val := range sb.slice {
		consumer.Accept(val)
	}
}

func (sb *spinedNodeBuilder[T]) AsSlice() []T {
	if sb.building {
		panic("spinedNodeBuilder: still building")
	}
	return sb.slice
}

func (sb *spinedNodeBuilder[T]) CopyInto(dst []T) {
	if sb.building {
		panic("spinedNodeBuilder: still building")
	}
	copy(dst, sb.slice)
}

func (sb *spinedNodeBuilder[T]) Count() int64 {
	return int64(len(sb.slice))
}
