package cgexec

import (
	"context"
	"time"
)

type IContext interface {
	Context() context.Context
	Done() <-chan struct{}
}

type Context struct {
	ctx context.Context
}

func NewContext(timer int) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timer)*time.Second)

	return Context{ctx}, cancel
}

func (c Context) Context() context.Context {
	return c.ctx
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}
