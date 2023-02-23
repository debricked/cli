package resolution

import (
	"errors"
	"github.com/debricked/cli/pkg/resolution/file"
	fileTestdata "github.com/debricked/cli/pkg/resolution/file/testdata"
	"github.com/debricked/cli/pkg/resolution/strategy"
	strategyTestdata "github.com/debricked/cli/pkg/resolution/strategy/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewResolver(t *testing.T) {
	r := NewResolver(
		file.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(10),
	)
	assert.NotNil(t, r)
}

func TestResolve(t *testing.T) {
	r := NewResolver(
		file.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(10),
	)

	res, err := r.Resolve([]string{"go.mod"})
	assert.NotEmpty(t, res.Jobs())
	assert.NoError(t, err)
}

func TestResolveStrategyError(t *testing.T) {
	r := NewResolver(
		fileTestdata.NewBatchFactoryMock(),
		strategy.NewStrategyFactory(),
		NewScheduler(10),
	)

	res, err := r.Resolve([]string{"go.mod"})
	assert.Empty(t, res.Jobs())
	assert.NoError(t, err)
}

func TestResolveScheduleError(t *testing.T) {
	errAssertion := errors.New("error")
	r := NewResolver(
		file.NewBatchFactory(),
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{Err: errAssertion},
	)

	res, err := r.Resolve([]string{"go.mod"})
	assert.NotEmpty(t, res.Jobs())
	assert.ErrorIs(t, err, errAssertion)
}
