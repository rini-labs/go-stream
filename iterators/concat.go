package iterators

import (
	"math"

	"github.com/rini-labs/go-stream"
)

func Concat[OUT any](a, b stream.Iterator[OUT]) stream.Iterator[OUT] {
	return &concatIterator[OUT]{
		a:       a,
		b:       b,
		unsized: a.GetExactSizeIfKnown() < 0 || b.GetExactSizeIfKnown() < 0,
	}
}

type concatIterator[OUT any] struct {
	a stream.Iterator[OUT]
	b stream.Iterator[OUT]

	processedA bool

	unsized bool
}

func (c *concatIterator[OUT]) Next() bool {
	hasNext := c.a.Next()
	if !hasNext {
		c.processedA = true
		hasNext = c.b.Next()
	}
	return hasNext
}

func (c *concatIterator[OUT]) Get() (OUT, error) {
	if c.processedA {
		return c.b.Get()
	}
	return c.a.Get()
}

func (c *concatIterator[OUT]) TryAdvance(consumer stream.Consumer[OUT]) bool {
	hasNext := c.a.TryAdvance(consumer)
	if !hasNext {
		c.processedA = true
		hasNext = c.b.TryAdvance(consumer)
	}
	return hasNext
}

func (c *concatIterator[OUT]) ForEachRemaining(consumer stream.Consumer[OUT]) {
	c.a.ForEachRemaining(consumer)
	c.processedA = true
	c.b.ForEachRemaining(consumer)
}

func (c *concatIterator[OUT]) EstimateSize() int {
	size := c.a.EstimateSize() + c.b.EstimateSize()
	if size >= 0 {
		return size
	}
	return math.MaxInt32
}

func (c *concatIterator[OUT]) GetExactSizeIfKnown() int {
	return getExactSizeIfKnown[OUT](c)
}

func (c *concatIterator[OUT]) Characteristics() int {
	var sizedFlag int
	if c.unsized {
		sizedFlag = stream.IsSized
	}

	return (c.a.Characteristics() & c.b.Characteristics()) ^ (stream.IsDistinct | stream.IsSorted | sizedFlag)
}

func (c *concatIterator[OUT]) HasCharacteristics(characteristics int) bool {
	return c.Characteristics()&characteristics == characteristics
}

var _ stream.Iterator[any] = &concatIterator[any]{}
