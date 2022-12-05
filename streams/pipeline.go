package streams

import (
	"errors"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/foreach"
	"github.com/rini-labs/go-stream/iterators"
	"github.com/rini-labs/go-stream/nodes"
	"github.com/rini-labs/go-stream/supplier"
)

type pipeline[OUT any] interface {
	Close() error
	markAsConsumed()
	iterator() (stream.Iterator[OUT], error)
}

func ofPipeline[IN any, OUT any](sourcePipeline pipeline[IN], wrapSink func(sink stream.Sink[OUT]) stream.Sink[IN]) *derivedPipeline[IN, OUT] {
	rv := &derivedPipeline[IN, OUT]{sourcePipeline: sourcePipeline, wrapSink: wrapSink}
	sourcePipeline.markAsConsumed()
	return rv
}

type derivedPipeline[IN any, OUT any] struct {
	linkedOrConsumed bool
	sourcePipeline   pipeline[IN]
	wrapSink         func(sink stream.Sink[OUT]) stream.Sink[IN]
}

// Close implements stream.Stream
func (p *derivedPipeline[IN, OUT]) Close() error {
	return p.sourcePipeline.Close()
}

func (p *derivedPipeline[IN, OUT]) markAsConsumed() {
	p.linkedOrConsumed = true
}

// Iterator implements stream.Stream
func (p *derivedPipeline[IN, OUT]) Iterator() (stream.Iterator[OUT], error) {
	if p.linkedOrConsumed {
		return nil, errors.New("stream has already been operated upon or closed")
	}

	return p.iterator()
}

func (p *derivedPipeline[IN, OUT]) iterator() (stream.Iterator[OUT], error) {
	return p.wrap(p.sourceIterator()), nil
}

func (p *derivedPipeline[IN, OUT]) sourceIterator() stream.Supplier[stream.Iterator[IN]] {
	return supplier.Of(p.sourcePipeline.iterator)
}

// ForEach implements stream.Stream
func (p *derivedPipeline[IN, OUT]) ForEach(consumer stream.Consumer[OUT]) error {
	iterator, err := p.Iterator()
	if err != nil {
		return err
	}
	// type args [OUT, any]
	evaluate(iterator, foreach.Of(consumer))
	return nil
}

func (p *derivedPipeline[IN, OUT]) Filter(predicate stream.Predicate[OUT]) stream.Stream[OUT] {
	return filterPipeline[OUT](p, predicate)
}

func (p *derivedPipeline[IN, OUT]) ToArray() ([]OUT, error) {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	return p.evaluateToArrayNode().AsArray()
}

func (p *derivedPipeline[IN, OUT]) wrap(iteratorSupplier stream.Supplier[stream.Iterator[IN]]) stream.Iterator[OUT] {
	return iterators.WrappingIterator(p.wrapSink, iteratorSupplier)
}

func (p *derivedPipeline[IN, OUT]) evaluateToArrayNode() nodes.Node[OUT] {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	iterator, err := p.sourceIterator().Get()
	if err != nil {
		panic(err)
	}
	return p.evaluate(iterator)
}

func (p *derivedPipeline[IN, OUT]) evaluate(iterator stream.Iterator[IN]) nodes.Node[OUT] {
	builder := nodes.BuilderOfSize[OUT](-1)
	return wrapAndcopyInto(builder, p.wrapSink, iterator).Build()
}

func wrapAndcopyInto[IN any, OUT any, S stream.Sink[OUT]](sink S, wrapSink func(sink stream.Sink[OUT]) stream.Sink[IN], iterator stream.Iterator[IN]) S {
	copyInto(wrapSink(sink), iterator)
	return sink
}
