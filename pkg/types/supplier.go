package types

type Supplier[T any] interface {
	Get() T
}

type supplier[T any] struct {
	f func() T
}

func (s *supplier[T]) Get() T {
	return s.f()
}

func NewSupplier[T any](f func() T) Supplier[T] {
	return &supplier[T]{f: f}
}
