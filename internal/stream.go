package internal

import (
	"github.com/rini-labs/go-stream/internal/iterator"
	"sort"

	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/internal/counting"
	"github.com/rini-labs/go-stream/internal/for_each"
	"github.com/rini-labs/go-stream/internal/reduce"
	"github.com/rini-labs/go-stream/pkg/types"
)

func NewStreamImpl[T any](iterator stream.Iterator[T]) stream.Stream[T] {
	return &streamImpl[T]{iterator: iterator}
}

// streamImpl is a generic stream that is iterated by the iterator returned by the
// supplier function
type streamImpl[T any] struct {
	iterator      stream.Iterator[T]
	closeHandlers []func() error
}

func (s *streamImpl[T]) OnClose(closeHandler func() error) stream.Stream[T] {
	s.closeHandlers = append(s.closeHandlers, closeHandler)
	return s
}

func (s *streamImpl[T]) Close() error {
	var err error
	for _, handler := range s.closeHandlers {
		if closeErr := handler(); err == nil && closeErr != nil {
			err = closeErr
		}
	}
	return err
}

func (s *streamImpl[T]) Iterator() stream.Iterator[T] {
	return s.iterator
}

func (s *streamImpl[T]) Filter(predicate types.Predicate[T]) stream.Stream[T] {
	var zeroValue T
	return NewStreamImpl(iterator.New(func() (T, error) {
		var nextValue T
		var err error
		for nextValue, err = s.iterator.Next(); err == nil; nextValue, err = s.iterator.Next() {
			if predicate(nextValue) {
				return nextValue, nil
			}
		}
		return zeroValue, err
	}))
}

func (s *streamImpl[T]) Map(mapper func(T) T) stream.Stream[T] {
	return Map[T, T](s, mapper)
}

func (s *streamImpl[T]) FlatMap(mapper func(T) stream.Stream[T]) stream.Stream[T] {
	return FlatMap[T, T](s, mapper)
}

func (s *streamImpl[T]) Peek(consumer func(T)) stream.Stream[T] {
	return NewStreamImpl(iterator.New(func() (T, error) {
		nextValue, err := s.iterator.Next()
		if err != nil {
			consumer(nextValue)
		}
		return nextValue, err
	}))
}

func (s *streamImpl[T]) Limit(maxSize int64) stream.Stream[T] {
	remaining := maxSize
	return NewStreamImpl(iterator.New(func() (T, error) {
		if remaining <= 0 {
			var zeroVal T
			return zeroVal, stream.Done
		}
		nextValue, err := s.iterator.Next()
		if err == nil {
			remaining--
		}
		return nextValue, err
	}))
}

func (s *streamImpl[T]) Sorted(comparator types.Comparator[T]) stream.Stream[T] {
	initialized := false
	var slice []T
	var zeroValue T
	return NewStreamImpl(iterator.New(func() (T, error) {
		if !initialized {
			slice = s.ToSlice()
			sort.Slice(slice, func(i, j int) bool {
				return comparator(slice[i], slice[j]) < 0
			})
			initialized = true
		}
		if len(slice) <= 0 {
			return zeroValue, stream.Done
		}
		rv := slice[0]
		slice = slice[1:]
		return rv, nil
	}))
}

func (s *streamImpl[T]) Skip(count int64) stream.Stream[T] {
	remaining := count
	return NewStreamImpl(iterator.New(func() (T, error) {
		for remaining > 0 {
			remaining--
			skippedVal, err := s.iterator.Next()
			if err != nil {
				return skippedVal, err
			}
		}
		return s.iterator.Next()
	}))
}

func (s *streamImpl[T]) ToSlice() []T {
	return evaluateToNode[T, T](s).AsSlice()
}

func (s *streamImpl[T]) ForEach(consumer types.Consumer[T]) {
	evaluate(s, for_each.MakeForEach(consumer))
}

func (s *streamImpl[T]) Count() int64 {
	return evaluate(s, counting.MakeCounting[T]())
}

func (s *streamImpl[T]) Reduce(accumulator types.BinaryOperator[T]) types.Optional[T] {
	return evaluate(s, reduce.MakeReducerFromOperator(accumulator))
}

func (s *streamImpl[T]) CopyInto(sink stream.Sink[T], iterator stream.Iterator[T]) {
	sink.Begin(0)
	iterator.ForEachRemaining(sink)
	sink.End()
}

func (s *streamImpl[T]) WrapAndCopyInto(sink stream.Sink[T], iterator stream.Iterator[T]) stream.Sink[T] {
	s.CopyInto(sink, iterator)
	return sink
}

func evaluate[IN any, OUT any](s *streamImpl[IN], terminalOp stream.TerminalOp[IN, OUT]) OUT {
	return terminalOp.EvaluateSequential(s, s.iterator)
}

func evaluateToNode[IN any, OUT any](s *streamImpl[IN]) Node[OUT] {
	nb := MakeNodeBuilder[OUT](-1)
	return wrapAndCopyInto[IN, OUT, NodeBuilder[OUT]](nb, s.iterator).Build()
}

func wrapAndCopyInto[IN any, OUT any, S stream.Sink[OUT]](sink S, iterator stream.Iterator[IN]) S {
	wrappedSink := wrapSink[IN, OUT](sink)
	CopyInto(wrappedSink, iterator)
	return sink
}

func CopyInto[IN any](wrappedSink stream.Sink[IN], i stream.Iterator[IN]) {
	wrappedSink.Begin(-1)
	i.ForEachRemaining(wrappedSink)
	wrappedSink.End()
}

// Needs to improve
func wrapSink[IN any, OUT any](s stream.Sink[OUT]) stream.Sink[IN] {
	return s.(stream.Sink[IN])
}
