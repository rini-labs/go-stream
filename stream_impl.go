package stream

type iterator[T any] struct {
	next func() (T, error)
}

func NewIterator[T any](next func() (T, error)) Iterator[T] {
	return &iterator[T]{next: next}
}

func (si *iterator[T]) Next() (T, error) {
	return si.next()
}

func NewStreamImpl[T any](iterator Iterator[T]) Stream[T] {
	return &streamImpl[T]{iterator: iterator}
}

// streamImpl is a generic stream that is iterated by the iterator returned by the
// supplier function
type streamImpl[T any] struct {
	iterator      Iterator[T]
	closeHandlers []func() error
}

func (s *streamImpl[T]) OnClose(closeHandler func() error) Stream[T] {
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

func (s *streamImpl[T]) Iterator() Iterator[T] {
	return s.iterator
}

func (s *streamImpl[T]) Filter(predicate Predicate[T]) Stream[T] {
	var zeroValue T
	return &streamImpl[T]{
		iterator: &iterator[T]{
			next: func() (T, error) {
				var nextValue T
				var err error
				for nextValue, err = s.iterator.Next(); err == nil; nextValue, err = s.iterator.Next() {
					if predicate(nextValue) {
						return nextValue, nil
					}
				}
				return zeroValue, err
			},
		},
	}
}

func (s *streamImpl[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](s, mapper)
}

func (s *streamImpl[T]) FlatMap(mapper func(T) Stream[T]) Stream[T] {
	return FlatMap[T, T](s, mapper)
}

func (s *streamImpl[T]) ToSlice() []T {
	var rv []T
	for nextVal, err := s.iterator.Next(); err == nil; nextVal, err = s.iterator.Next() {
		rv = append(rv, nextVal)
	}
	return rv
}
