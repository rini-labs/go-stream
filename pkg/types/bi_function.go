package types

type BiFunction[T any, U any, R any] interface {
	Apply(t T, u U) R
}

func BiFunctionFromFunc[T any, U any, R any](f func(T, U) R) BiFunction[T, U, R] {
	return &biFunction[T, U, R]{f: f}
}

type biFunction[T any, U any, R any] struct {
	f func(T, U) R
}

func (bf *biFunction[T, U, R]) Apply(t T, u U) R {
	return bf.f(t, u)
}

type BinaryOperator[T any] interface {
	BiFunction[T, T, T]
}

func BinaryOperatorFromFunc[T any](f func(T, T) T) BinaryOperator[T] {
	return &biFunction[T, T, T]{f: f}
}
