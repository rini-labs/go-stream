package stream

type Iterator[OUT any] interface {
	Next() bool
	Get() (OUT, error)
	TryAdvance(consumer Consumer[OUT]) bool
	ForEachRemaining(consumer Consumer[OUT])
}

type Iterable[OUT any] interface {
	Iterator() (Iterator[OUT], error)
	ForEach(consumer Consumer[OUT]) error
}
