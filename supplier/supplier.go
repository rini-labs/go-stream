package supplier

import "github.com/rini-labs/go-stream"

func Of[T any](provider func() (T, error)) stream.Supplier[T] {
	return &supplier[T]{provider: provider}
}

type supplier[T any] struct {
	provider func() (T, error)
}

func (s *supplier[T]) Get() (T, error) {
	return s.provider()
}
