package iterators

import (
	"github.com/rini-labs/go-stream"
)

type PipelineHelper[IN any, OUT any] interface {
	GetStreamAndOpFlags() int
	WrapSink(sink stream.Sink[OUT]) stream.Sink[IN]
}

func WrappingIterator[IN any, OUT any](ph PipelineHelper[IN, OUT], iteratorSupplier stream.Supplier[stream.Iterator[IN]]) stream.Iterator[OUT] {
	return &wrappingIterator[IN, OUT]{
		ph:               ph,
		iteratorSupplier: iteratorSupplier,
	}
}

type wrappingIterator[IN any, OUT any] struct {
	ph               PipelineHelper[IN, OUT]
	iteratorSupplier stream.Supplier[stream.Iterator[IN]]

	iterator   stream.Iterator[IN]
	bufferSink stream.Sink[IN]
	buffer     Buffer[OUT]

	nextToConsume int
	finished      bool
}

// Next implements stream.Iterator
func (wi *wrappingIterator[IN, OUT]) Next() bool {
	if wi.buffer == nil {
		if wi.finished {
			return false
		}

		wi.init()
		wi.nextToConsume = 0
		wi.bufferSink.Begin(-1)
		return wi.fillBuffer()
	}
	wi.nextToConsume++
	hasNext := wi.nextToConsume < wi.buffer.Count()
	if !hasNext {
		wi.nextToConsume = 0
		wi.buffer.Clear()
		hasNext = wi.fillBuffer()
	}
	return hasNext
}

// Get implements stream.Iterator
func (wi *wrappingIterator[IN, OUT]) Get() (OUT, error) {
	return wi.buffer.Get(wi.nextToConsume)
}

// TryAdvance implements stream.Iterator
func (wi *wrappingIterator[IN, OUT]) TryAdvance(consumer stream.Consumer[OUT]) bool {
	return tryAdvance[OUT](wi, consumer)
}

// ForEachRemaining implements stream.Iterator
func (wi *wrappingIterator[IN, OUT]) ForEachRemaining(consumer stream.Consumer[OUT]) {
	forEachRemaining[OUT](wi, consumer)
}

func (wi *wrappingIterator[IN, OUT]) EstimateSize() int {
	wi.init()
	return wi.iterator.EstimateSize()
}
func (wi *wrappingIterator[IN, OUT]) GetExactSizeIfKnown() int {
	if stream.OpFlagSized.IsKnown(wi.ph.GetStreamAndOpFlags()) {
		return wi.iterator.GetExactSizeIfKnown()
	}
	return -1
}
func (wi *wrappingIterator[IN, OUT]) Characteristics() int {
	wi.init()

	c := stream.ToCharacteristics(stream.ToStreamFlags(wi.ph.GetStreamAndOpFlags()))
	if (c & SIZED) != 0 {
		c &= ^SIZED
		c |= wi.iterator.Characteristics() & SIZED
	}

	return c
}
func (wi *wrappingIterator[IN, OUT]) HasCharacteristics(characteristics int) bool {
	return hasCharacteristics[OUT](wi, characteristics)
}

func (wi *wrappingIterator[IN, OUT]) init() error {
	var err error
	// Convert iteratorSupplier to iterator
	wi.iterator, err = wi.iteratorSupplier.Get()
	if err != nil {
		return err
	}

	b := NewBuffer[OUT]()
	wi.buffer = b
	wi.bufferSink = wi.ph.WrapSink(b)
	return nil
}

// if the buffer is empty, push elements into the sink chain until the source is empty or cancellation is requested.
func (wi *wrappingIterator[IN, OUT]) fillBuffer() bool {
	for wi.buffer.Count() == 0 {
		// Try consuming from the source iterator
		if !wi.iterator.TryAdvance(wi.bufferSink) {
			// If the source iterator is empty, we're done
			if wi.finished {
				return false
			} else {
				wi.bufferSink.End()
				wi.finished = true
			}
		}
	}
	return true
}
