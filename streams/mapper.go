package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/mapper"
)

func mapPipeline[IN any, OUT any](p pipeline[IN], mapperFunc func(v IN) OUT) stream.Stream[OUT] {
	return OfStateless(p, func(flags int, sink stream.Sink[OUT]) stream.Sink[IN] {
		return mapper.NewMapSink(sink, mapperFunc)
	}, stream.NotSorted|stream.NotDistinct)
}

func flatMapPipeline[IN any, OUT any](p pipeline[IN], mapperFunc func(v IN) stream.Stream[OUT]) stream.Stream[OUT] {
	return OfStateless(p, func(flags int, sink stream.Sink[OUT]) stream.Sink[IN] {
		return mapper.NewFlatMapSink(sink, mapperFunc)
	}, stream.NotSorted|stream.NotDistinct|stream.NotSized)
}

func Map[IN any, OUT any](s stream.Stream[IN], mapper func(v IN) OUT) stream.Stream[OUT] {
	return mapPipeline(s.(pipeline[IN]), mapper)
}

func FlatMap[IN any, OUT any](s stream.Stream[IN], mapper func(v IN) stream.Stream[OUT]) stream.Stream[OUT] {
	return flatMapPipeline(s.(pipeline[IN]), mapper)
}
