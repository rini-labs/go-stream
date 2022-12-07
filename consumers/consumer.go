package consumers

import "github.com/rini-labs/go-stream"

func Of[IN any](f func(val IN)) stream.Consumer[IN] {
	return &consumer[IN]{f: f}
}

type consumer[IN any] struct {
	f func(val IN)
}

func (c *consumer[IN]) Accept(val IN) {
	c.f(val)
}
