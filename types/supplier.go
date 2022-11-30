package types

type Supplier[T any] interface {
	Get() T
}
