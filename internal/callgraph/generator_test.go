package callgraph

import (
	"errors"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/debricked/cli/internal/callgraph/config"
	strategyTestdata "github.com/debricked/cli/internal/callgraph/strategy/testdata"
	"github.com/debricked/cli/internal/io/finder/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	workers   = 10
	goModFile = "go.mod"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator(
		&testdata.FinderMock{},
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)
	assert.NotNil(t, g)
}

func TestGenerate(t *testing.T) {
	g := NewGenerator(
		&testdata.FinderMock{},
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate([]string{"../../go.mod"}, nil, configs, ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, g.Generation.Jobs())
}

func TestGenerateWithTimer(t *testing.T) {
	g := NewGenerator(
		&testdata.FinderMock{},
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	err := g.GenerateWithTimer([]string{"../../go.mod"}, nil, configs, 1000)
	assert.NoError(t, err)
	assert.NotEmpty(t, g.Generation.Jobs())
}

func TestGenerateInvokeError(t *testing.T) {
	g := NewGenerator(
		&testdata.FinderMock{},
		strategyTestdata.NewStrategyFactoryErrorMock(),
		NewScheduler(workers),
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate([]string{"../../go.mod"}, nil, configs, ctx)
	assert.NotNil(t, err)
}

func TestGenerateScheduleError(t *testing.T) {
	errAssertion := errors.New("error")
	g := NewGenerator(
		&testdata.FinderMock{},
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{Err: errAssertion},
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate([]string{"../../go.mod"}, nil, configs, ctx)
	assert.NotEmpty(t, g.Generation.Jobs())
	assert.ErrorIs(t, err, errAssertion)
}

func TestGenerateDirWithoutConfig(t *testing.T) {
	g := NewGenerator(
		&testdata.FinderMock{},
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate([]string{"invalid-dir"}, nil, nil, ctx)
	assert.Empty(t, g.Generation.Jobs())
	assert.NoError(t, err)
}
