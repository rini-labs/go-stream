package types

type Consumer[T any] interface {
	Accept(val T)
}

type consumerFunc[T any] struct {
	fn func(val T)
}

func (c consumerFunc[T]) Accept(val T) {
	c.fn(val)
}

func ToConsumer[T any](f func(val T)) Consumer[T] {
	return consumerFunc[T]{fn: f}
}
