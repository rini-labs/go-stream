package stream

type Iterator[OUT any] interface {
	Next() bool
	Get() (OUT, error)
	TryAdvance(consumer Consumer[OUT]) bool
	ForEachRemaining(consumer Consumer[OUT])

	EstimateSize() int
	GetExactSizeIfKnown() int
	Characteristics() int
	HasCharacteristics(characteristics int) bool
}

type Iterable[OUT any] interface {
	Iterator() (Iterator[OUT], error)
	ForEach(consumer Consumer[OUT]) error
}
