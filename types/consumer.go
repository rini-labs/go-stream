package types

type Consumer[T any] interface {
	Accept(val T)
}
