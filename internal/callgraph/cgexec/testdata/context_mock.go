package testdata

import (
	"context"
	"time"
)

type ContextMock struct {
	ctx context.Context
}

func NewContextMock() (ContextMock, context.CancelFunc) {
	ctx := context.Background()
	return ContextMock{ctx}, nil
}

func NewContextMockCancelled() (ContextMock, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return ContextMock{ctx}, nil
}

func NewContextMockDeadlineReached() (ContextMock, context.CancelFunc) {
	ctx := context.Background()
	ctx, _ = context.WithDeadline(ctx, time.Now())
	return ContextMock{ctx}, nil
}

func (c ContextMock) Context() context.Context {
	return c.ctx
}

func (c ContextMock) Done() <-chan struct{} {
	context.WithCancel(context.Background())
	return c.ctx.Done()
}

func (c ContextMock) Err() error {
	return context.DeadlineExceeded
}
