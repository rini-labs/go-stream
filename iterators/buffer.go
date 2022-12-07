package iterators

import (
	"errors"

	"github.com/rini-labs/go-stream"
)

type Buffer[OUT any] interface {
	stream.Sink[OUT]
	stream.Iterable[OUT]

	Count() int
	Clear()
	Get(index int) (OUT, error)
	CopyInto([]OUT)
}

var (
	_ Buffer[any] = (*buffer[any])(nil)
)

func NewBuffer[OUT any]() Buffer[OUT] {
	return &buffer[OUT]{}
}

type buffer[OUT any] struct {
	buffer []OUT
}

// Begin implements Buffer
func (b *buffer[OUT]) Begin(size int) {
	if size != -1 {
		b.buffer = make([]OUT, 0, size)
	} else {
		b.buffer = make([]OUT, 0)
	}
}

// End implements Buffer
func (b *buffer[OUT]) End() {
}

// Accept implements Buffer
func (b *buffer[OUT]) Accept(val OUT) {
	b.buffer = append(b.buffer, val)
}

// Iterator implements Buffer
func (b *buffer[OUT]) Iterator() (stream.Iterator[OUT], error) {
	return OfSlice(b.buffer), nil
}

func (b *buffer[OUT]) ForEach(consumer stream.Consumer[OUT]) error {
	for _, val := range b.buffer {
		consumer.Accept(val)
	}
	return nil
}

// Count implements Buffer
func (b *buffer[OUT]) Count() int {
	return len(b.buffer)
}

// Count implements Buffer
func (b *buffer[OUT]) Clear() {
	b.buffer = nil
}

func (b *buffer[OUT]) Get(index int) (OUT, error) {
	var zeroValue OUT
	if index < 0 || index >= len(b.buffer) {
		return zeroValue, errors.New("index out of range")
	}
	return b.buffer[index], nil
}

func (b *buffer[OUT]) CopyInto(data []OUT) {
	copy(data, b.buffer)
}
