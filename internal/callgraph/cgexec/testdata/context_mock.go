package testdata

import "context"

type ContextMock struct {
	ctx context.Context
}

func NewContextMock() (ContextMock, context.CancelFunc) {
	ctx := context.Background()
	return ContextMock{ctx}, nil
}

func (c ContextMock) Context() context.Context {
	return c.ctx
}

func (c ContextMock) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c ContextMock) Err() error {
	return context.DeadlineExceeded
}
