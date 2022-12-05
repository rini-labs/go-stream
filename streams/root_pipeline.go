package streams

import (
	"errors"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/foreach"
	"github.com/rini-labs/go-stream/nodes"
	"github.com/rini-labs/go-stream/sinks"
	"github.com/rini-labs/go-stream/supplier"
)

func Of[OUT any](iterator stream.Iterator[OUT]) stream.Stream[OUT] {
	iteratorSupplier := supplier.Of(func() (stream.Iterator[OUT], error) { return iterator, nil })
	return OfSupplier(iteratorSupplier)
}

func OfSupplier[OUT any](iteratorSupplier stream.Supplier[stream.Iterator[OUT]]) stream.Stream[OUT] {
	return &rootPipeline[OUT]{sourceIteratorSupplier: iteratorSupplier}
}

type rootPipeline[OUT any] struct {
	sourceIteratorSupplier stream.Supplier[stream.Iterator[OUT]]
	linkedOrConsumed       bool
}

// Close implements Stream
func (p *rootPipeline[OUT]) Close() error {
	return nil
}

func (p *rootPipeline[OUT]) markAsConsumed() {
	p.linkedOrConsumed = true
}

// Iterator implements Stream
func (p *rootPipeline[OUT]) Iterator() (stream.Iterator[OUT], error) {
	if p.linkedOrConsumed {
		return nil, errors.New("stream has already been operated upon or closed")
	}
	p.linkedOrConsumed = true
	return p.iterator()
}

func (p *rootPipeline[OUT]) iterator() (stream.Iterator[OUT], error) {
	if p.sourceIteratorSupplier != nil {
		iterator, err := p.sourceIteratorSupplier.Get()
		if err != nil {
			return nil, err
		}
		// Clear the sourceIteratorSupplier after consuming it
		p.sourceIteratorSupplier = nil
		return iterator.(stream.Iterator[OUT]), nil
	}
	return nil, errors.New("stream already consumed")
}

func (p *rootPipeline[OUT]) ForEach(consumer stream.Consumer[OUT]) error {
	iterator, err := p.Iterator()
	if err != nil {
		return err
	}
	// type args [OUT, any]
	evaluate(iterator, foreach.Of(consumer))
	return nil
}

func (p *rootPipeline[OUT]) Filter(predicate stream.Predicate[OUT]) stream.Stream[OUT] {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	return filterPipeline[OUT](p, predicate)
}

func (p *rootPipeline[OUT]) ToArray() ([]OUT, error) {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	return p.evaluateToArrayNode().AsArray()
}

func (p *rootPipeline[OUT]) evaluateToArrayNode() nodes.Node[OUT] {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	iterator, err := p.iterator()
	if err != nil {
		panic(err)
	}
	return p.evaluate(iterator)
}

func (p *rootPipeline[OUT]) evaluate(iterator stream.Iterator[OUT]) nodes.Node[OUT] {
	builder := nodes.BuilderOfSize[OUT](-1)
	return copyInto(builder, iterator).Build()
}

func evaluate[E_IN any, E_OUT any](iterator stream.Iterator[E_IN], terminalOp stream.TerminalOp[E_IN, E_OUT]) E_OUT {
	return terminalOp.EvaluateSequential(iterator)
}

func copyInto[OUT any, S stream.Sink[OUT]](sink S, iterator stream.Iterator[OUT]) S {
	sinks.CopyInto[OUT](iterator, sink)
	return sink
}
