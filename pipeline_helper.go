package stream

type PipelineHelper[IN any] interface {
	WrapAndCopyInto(sink Sink[IN], iterator Iterator[IN]) Sink[IN]
	CopyInto(sink Sink[IN], iterator Iterator[IN])
}
