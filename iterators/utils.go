package iterators

import "github.com/rini-labs/go-stream"

func tryAdvance[OUT any](iterator stream.Iterator[OUT], consumer stream.Consumer[OUT]) bool {
	if !iterator.Next() {
		return false
	}
	value, err := iterator.Get()
	if err != nil {
		return false
	}
	consumer.Accept(value)
	return true
}

func forEachRemaining[OUT any](iterator stream.Iterator[OUT], consumer stream.Consumer[OUT]) {
	for iterator.TryAdvance(consumer) {
	}
}

func getExactSizeIfKnown[OUT any](i stream.Iterator[OUT]) int {
	if i.HasCharacteristics(SIZED) {
		return -1
	}
	return i.EstimateSize()
}

func hasCharacteristics[OUT any](i stream.Iterator[OUT], characteristics int) bool {
	return (i.Characteristics() & characteristics) == characteristics
}
