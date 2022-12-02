package types

import "errors"

// Optional represents value that may be present or not.
type Optional[T any] interface {
	// IsPresent returns true if value is present.
	IsPresent() bool

	// Get returns value if present, otherwise returns error.
	Get() (T, error)

	// OrElse returns value if present, otherwise returns other.
	OrElse(other T) T

	// OrElseGet returns value if present, otherwise calls supplier and returns its result.
	OrElseGet(supplier Supplier[T]) T

	// OrElseReturnError returns value if present, otherwise return the error returned by the supplier.
	OrElseReturnError(supplier Supplier[error]) (T, error)

	// IfPresent calls consumer with value if present.
	IfPresent(consumer Consumer[T])

	// Filter returns an Optional describing the value if it is present and matches the given predicate.
	Filter(predicate Predicate[T]) Optional[T]
}

type optionImpl[T any] struct {
	value T
	err   error
}

// OptionalOf returns an Optional describing the specified value, if non-null, otherwise returns an empty Optional.
func OptionalOf[T any](value T) Optional[T] {
	return &optionImpl[T]{value: value}
}

func EmptyOptional[T any]() Optional[T] {
	return &optionImpl[T]{err: errors.New("empty")}
}

func (o *optionImpl[T]) IsPresent() bool {
	return o.err == nil
}

func (o *optionImpl[T]) Get() (T, error) {
	return o.value, o.err
}

func (o *optionImpl[T]) OrElse(other T) T {
	if o.err != nil {
		return other
	}
	return o.value
}

func (o *optionImpl[T]) OrElseGet(supplier Supplier[T]) T {
	if o.err != nil {
		return supplier.Get()
	}
	return o.value
}

func (o *optionImpl[T]) OrElseReturnError(supplier Supplier[error]) (T, error) {
	if o.err != nil {
		return o.value, supplier.Get()
	}
	return o.value, nil
}

func (o *optionImpl[T]) IfPresent(consumer Consumer[T]) {
	if o.err == nil {
		consumer.Accept(o.value)
	}
}

func (o *optionImpl[T]) Filter(predicate Predicate[T]) Optional[T] {
	if o.err != nil || predicate(o.value) {
		return o
	}
	return EmptyOptional[T]()
}
