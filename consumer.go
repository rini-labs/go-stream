package stream

type Consumer[IN any] interface {
	Accept(IN)
}
