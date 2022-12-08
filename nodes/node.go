package nodes

import "github.com/rini-labs/go-stream"

type Node[T any] interface {
	stream.Iterable[T]

	Count() int

	AsSlice() ([]T, error)
}

type NodeBuilder[T any] interface {
	stream.Sink[T]

	Build() Node[T]
}

func BuilderOfSize[T any](size int) NodeBuilder[T] {
	if size != -1 {
		return NewSliceBuilder[T](size)
	}
	return Builder[T]()
}
