package stream

type Stream[OUT any] interface {
	Close() error

	Iterable[OUT]

	Filter(predicate Predicate[OUT]) Stream[OUT]

	ToArray() ([]OUT, error)

	Sort(comparator Comparator[OUT]) Stream[OUT]

	Peek(consumer Consumer[OUT]) Stream[OUT]

	Limit(limit int) Stream[OUT]

	Skip(skip int) Stream[OUT]
}
