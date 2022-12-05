package streams

import (
	"github.com/rini-labs/go-stream"
	"github.com/rini-labs/go-stream/mapper"
)

func mapPipeline[IN any, OUT any](p pipeline[IN], mapperFunc func(v IN) OUT) *derivedPipeline[IN, OUT] {
	return ofPipeline(p, func(sink stream.Sink[OUT]) stream.Sink[IN] {
		return mapper.NewMapSink(sink, mapperFunc)
	})
}

func flatMapPipeline[IN any, OUT any](p pipeline[IN], mapperFunc func(v IN) stream.Stream[OUT]) *derivedPipeline[IN, OUT] {
	return ofPipeline(p, func(sink stream.Sink[OUT]) stream.Sink[IN] {
		return mapper.NewFlatMapSink(sink, mapperFunc)
	})
}

func Map[IN any, OUT any](s stream.Stream[IN], mapper func(v IN) OUT) stream.Stream[OUT] {
	return mapPipeline(s.(pipeline[IN]), mapper)
}

func FlatMap[IN any, OUT any](s stream.Stream[IN], mapper func(v IN) stream.Stream[OUT]) stream.Stream[OUT] {
	return flatMapPipeline(s.(pipeline[IN]), mapper)
}
