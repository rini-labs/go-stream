package stream

type Comparator[T any] interface {
	Compare(val1 T, val2 T) int
}
