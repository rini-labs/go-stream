package types

type BiFunction[T any, U any, R any] interface {
	Apply(t T, u U) R
}
