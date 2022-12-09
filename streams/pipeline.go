package streams

import (
	"errors"
	"github.com/rini-labs/go-stream/iterators"
	"github.com/rini-labs/go-stream/reduce"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/foreach"
	"github.com/rini-labs/go-stream/nodes"
	"github.com/rini-labs/go-stream/sinks"
	"github.com/rini-labs/go-stream/supplier"
)

func Of[OUT any](iterator stream.Iterator[OUT], sourceFlags int) stream.Stream[OUT] {
	iteratorSupplier := supplier.Of(func() (stream.Iterator[OUT], error) { return iterator, nil })
	return OfSupplier(iteratorSupplier, sourceFlags)
}

func OfSupplier[OUT any](iteratorSupplier stream.Supplier[stream.Iterator[OUT]], sourceFlags int) stream.Stream[OUT] {
	sourceOrOpFlags := sourceFlags & stream.MaskStream
	return &rootPipeline[OUT, OUT]{
		iteratorSupplier: iteratorSupplier,
		wrapSink:         nil,
		sourceOrOpFlags:  sourceOrOpFlags,
		combinedFlags:    ^(sourceOrOpFlags << 1) & stream.InitialOpsValue,
	}
}

type pipeline[OUT any] interface {
	Close() error
	markAsConsumed()
	iterator() (stream.Iterator[OUT], error)
	getCombinedFlags() int
}

type OpWrapSink[IN any, OUT any] func(flags int, sink stream.Sink[OUT]) stream.Sink[IN]

func OfStateful[IN any, OUT any](sourcePipeline pipeline[IN], wrapSink OpWrapSink[IN, OUT], opsFlag int) stream.Stream[OUT] {
	return ofPipeline(sourcePipeline, wrapSink, true, opsFlag)
}

func OfStateless[IN any, OUT any](sourcePipeline pipeline[IN], wrapSink OpWrapSink[IN, OUT], opsFlag int) stream.Stream[OUT] {
	return ofPipeline(sourcePipeline, wrapSink, false, opsFlag)
}

func ofPipeline[IN any, OUT any](sourcePipeline pipeline[IN], wrapSink OpWrapSink[IN, OUT], stateful bool, opsFlag int) stream.Stream[OUT] {
	rv := &rootPipeline[IN, OUT]{
		iteratorSupplier: supplier.Of(sourcePipeline.iterator),
		wrapSink:         wrapSink,
		sourcePipeline:   sourcePipeline,
		sourceOrOpFlags:  opsFlag & stream.MaskOp,
		combinedFlags:    stream.CombineOpFlags(opsFlag, sourcePipeline.getCombinedFlags()),
		opIsStateful:     stateful,
	}
	sourcePipeline.markAsConsumed()
	return rv
}

type rootPipeline[IN any, OUT any] struct {
	linkedOrConsumed bool
	iteratorSupplier stream.Supplier[stream.Iterator[IN]]
	wrapSink         OpWrapSink[IN, OUT]
	sourcePipeline   pipeline[IN]
	sourceOrOpFlags  int
	combinedFlags    int
	opIsStateful     bool
}

func (p *rootPipeline[IN, OUT]) getCombinedFlags() int {
	return p.combinedFlags
}

func (p *rootPipeline[IN, OUT]) Close() error {
	if p.sourcePipeline != nil {
		return p.sourcePipeline.Close()
	}
	return nil
}

func (p *rootPipeline[IN, OUT]) markAsConsumed() {
	p.linkedOrConsumed = true
}

// Iterator implements Stream
func (p *rootPipeline[IN, OUT]) Iterator() (stream.Iterator[OUT], error) {
	if p.linkedOrConsumed {
		return nil, errors.New("stream has already been operated upon or closed")
	}
	p.linkedOrConsumed = true
	return p.iterator()
}

func (p *rootPipeline[IN, OUT]) iterator() (stream.Iterator[OUT], error) {
	if p.wrapSink == nil {
		iter, err := p.iteratorSupplier.Get()
		if err != nil {
			return nil, err
		}
		return iter.(stream.Iterator[OUT]), nil
	}
	return iterators.WrappingIterator[IN, OUT](p, p.iteratorSupplier), nil
}

func (p *rootPipeline[IN, OUT]) sourceIterator() (stream.Iterator[IN], error) {
	return p.iteratorSupplier.Get()
}

// ForEach implements stream.Stream
func (p *rootPipeline[IN, OUT]) ForEach(consumer stream.Consumer[OUT]) error {

	iterator, err := p.Iterator()
	if err != nil {
		return err
	}
	// type args [OUT, any]
	evaluate(iterator, foreach.Of(consumer))
	return nil
}

func (p *rootPipeline[IN, OUT]) Filter(predicate stream.Predicate[OUT]) stream.Stream[OUT] {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	return filterPipeline[OUT](p, predicate)
}

func (p *rootPipeline[IN, OUT]) Sort(comparator stream.Comparator[OUT]) stream.Stream[OUT] {
	return sortPipeline[OUT](p, comparator, false)
}

func (p *rootPipeline[IN, OUT]) Peek(consumer stream.Consumer[OUT]) stream.Stream[OUT] {
	return peekPipeline[OUT](p, consumer)
}

func (p *rootPipeline[IN, OUT]) Limit(limit int) stream.Stream[OUT] {
	return slicePipeline[OUT](p, 0, limit)
}

func (p *rootPipeline[IN, OUT]) Skip(skip int) stream.Stream[OUT] {
	return slicePipeline[OUT](p, skip, -1)
}

func (p *rootPipeline[IN, OUT]) ReduceWithSeed(seed OUT, reducer func(OUT, OUT) OUT) (OUT, error) {
	iterator, err := p.Iterator()
	if err != nil {
		var zeroVal OUT
		return zeroVal, err
	}
	return evaluate[OUT, OUT](iterator, reduce.NewOpWithSeed(seed, reducer)), nil
}

func (p *rootPipeline[IN, OUT]) Reduce(reducer func(OUT, OUT) OUT) (OUT, error) {
	iterator, err := p.Iterator()
	if err != nil {
		var zeroVal OUT
		return zeroVal, err
	}
	return evaluate[OUT, OUT](iterator, reduce.NewOp(reducer)), nil
}

func (p *rootPipeline[IN, OUT]) Count() (int, error) {
	iterator, err := p.Iterator()
	if err != nil {
		return 0, err
	}
	return evaluate[OUT, int](iterator, reduce.NewOp(func(state int, value OUT) int {
		return state + 1
	})), nil
}

func (p *rootPipeline[IN, OUT]) ToSlice() ([]OUT, error) {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	return p.evaluateToArrayNode().AsSlice()
}

func (p *rootPipeline[IN, OUT]) WrapSink(sink stream.Sink[OUT]) stream.Sink[IN] {
	if p.wrapSink == nil {
		return sink.(stream.Sink[IN])
	}
	return p.wrapSink(p.getCombinedFlags(), sink)
}

func (p *rootPipeline[IN, OUT]) GetStreamAndOpFlags() int {
	return p.combinedFlags
}

func (p *rootPipeline[IN, OUT]) evaluateToArrayNode() nodes.Node[OUT] {
	if p.linkedOrConsumed {
		panic("stream has already been operated upon or closed")
	}
	iterator, err := p.sourceIterator()
	if err != nil {
		panic(err)
	}
	return p.evaluate(iterator)
}

func (p *rootPipeline[IN, OUT]) evaluate(iterator stream.Iterator[IN]) nodes.Node[OUT] {
	builder := nodes.BuilderOfSize[OUT](iterator.GetExactSizeIfKnown())
	return wrapAndCopyInto(builder, p.wrapSink, iterator).Build()
}

func evaluate[IN any, OUT any](iterator stream.Iterator[IN], terminalOp stream.TerminalOp[IN, OUT]) OUT {
	return terminalOp.EvaluateSequential(iterator)
}

func wrapAndCopyInto[IN any, OUT any, S stream.Sink[OUT]](sink S, wrapSink func(flags int, sink stream.Sink[OUT]) stream.Sink[IN], iterator stream.Iterator[IN]) S {
	copyInto(wrapSink(0, sink), iterator)
	return sink
}

func copyInto[OUT any, S stream.Sink[OUT]](sink S, iterator stream.Iterator[OUT]) S {
	sinks.CopyInto[OUT](iterator, sink)
	return sink
}
