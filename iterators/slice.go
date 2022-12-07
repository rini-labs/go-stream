package iterators

import "github.com/rini-labs/go-stream"

func OfSlice[OUT any](values []OUT) stream.Iterator[OUT] {
	return &sliceIterator[OUT]{values: values, index: -1}
}

type sliceIterator[OUT any] struct {
	values []OUT
	index  int
}

// Get implements stream.Iterator
func (si *sliceIterator[OUT]) Get() (OUT, error) {
	var zeroValue OUT
	if si.index < 0 {
		return zeroValue, stream.ErrIteratorNotStarted
	}
	if si.index >= len(si.values) {
		return zeroValue, stream.ErrDone
	}

	return si.values[si.index], nil
}

// Next implements stream.Iterator
func (si *sliceIterator[OUT]) Next() bool {
	si.index++
	if si.index >= len(si.values) {
		return false
	}
	return true
}

// ForEachRemaining implements stream.Iterator
func (si *sliceIterator[OUT]) ForEachRemaining(consumer stream.Consumer[OUT]) {
	forEachRemaining[OUT](si, consumer)
}

// TryAdvance implements stream.Iterator
func (si *sliceIterator[OUT]) TryAdvance(consumer stream.Consumer[OUT]) bool {
	return tryAdvance[OUT](si, consumer)
}

func (si *sliceIterator[OUT]) EstimateSize() int {
	return len(si.values)
}
func (si *sliceIterator[OUT]) GetExactSizeIfKnown() int {
	return getExactSizeIfKnown[OUT](si)
}
func (si *sliceIterator[OUT]) Characteristics() int {
	return ORDERED | SIZED
}
func (si *sliceIterator[OUT]) HasCharacteristics(characteristics int) bool {
	return hasCharacteristics[OUT](si, characteristics)
}
