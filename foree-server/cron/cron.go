package cron

import "context"

type Cron[T any] struct {
	ID               string
	IntervalInSecond int
	Receiver         chan T
	Initializer      func(context.Context)
	Handler          func(context.Context, T)
}

func (c *Cron[T]) Initial() {

}

func (c *Cron[T]) Submit(t T) {
	c.Receiver <- t
}
