package stream

type Sink[IN any] interface {
	Consumer[IN]

	Begin(size int)

	End()
}
