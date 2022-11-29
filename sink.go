package stream

type Consumer[T any] interface {
	Accept(val T)
}

type Sink[T any] interface {
	Consumer[T]

	// Begin resets the sink state to receive a fresh data set.
	Begin(size int64)

	// End indicates that all elements have been pushed.
	End()
}

type abstractSink[T any] struct {
}

func (s *abstractSink[T]) Begin(_ int64) {
}

func (s *abstractSink[T]) End() {
}

type chainedSink[T any] struct {
	abstractSink[T]
	downstream Sink[T]
}

func (s *chainedSink[T]) Begin(size int64) {
	s.downstream.Begin(size)
}

func (s *chainedSink[T]) End() {
	s.downstream.End()
}
