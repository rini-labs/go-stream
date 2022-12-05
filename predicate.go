package stream

type Predicate[IN any] interface {
	Test(value IN) bool
}
