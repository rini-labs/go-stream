package stream

type Supplier[T any] interface {
	Get() (T, error)
}
