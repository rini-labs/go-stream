package stream

type Stream[OUT any] interface {
	Close() error

	Iterable[OUT]

	Filter(predicate Predicate[OUT]) Stream[OUT]

	ToArray() ([]OUT, error)
}
