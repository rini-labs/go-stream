package predicates

import "github.com/rini-labs/go-stream"

func Of[T any](test func(value T) bool) stream.Predicate[T] {
	return &predicate[T]{test: test}
}

type predicate[T any] struct {
	test func(value T) bool
}

func (p *predicate[T]) Test(value T) bool {
	return p.test(value)
}
