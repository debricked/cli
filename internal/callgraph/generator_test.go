package callgraph

import (
	"errors"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/debricked/cli/internal/callgraph/config"
	strategyTestdata "github.com/debricked/cli/internal/callgraph/strategy/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	workers   = 10
	goModFile = "go.mod"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)
	assert.NotNil(t, g)
}

func TestGenerate(t *testing.T) {
	g := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}, true, "maven", ""),
	}
	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate(
		DebrickedOptions{
			Paths:      []string{"../../go.mod"},
			Exclusions: []string{},
			Inclusions: []string{},
			Configs:    configs,
		}, ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, g.Generation.Jobs())
}

func TestGenerateWithTimer(t *testing.T) {
	g := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}, true, "maven", ""),
	}
	err := g.GenerateWithTimer(
		DebrickedOptions{
			Paths:      []string{"../../go.mod"},
			Exclusions: nil,
			Inclusions: nil,
			Configs:    configs,
			Timeout:    1000,
		},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, g.Generation.Jobs())
}

func TestGenerateInvokeError(t *testing.T) {
	g := NewGenerator(
		strategyTestdata.NewStrategyFactoryErrorMock(),
		NewScheduler(workers),
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}, true, "maven", ""),
	}
	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate(
		DebrickedOptions{
			Paths:      []string{"../../go.mod"},
			Exclusions: []string{},
			Inclusions: []string{},
			Configs:    configs,
		}, ctx)
	assert.NotNil(t, err)
}

func TestGenerateScheduleError(t *testing.T) {
	errAssertion := errors.New("error")
	g := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{Err: errAssertion},
	)

	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}, true, "maven", ""),
	}
	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate(
		DebrickedOptions{
			Paths:      []string{"../../go.mod"},
			Exclusions: []string{},
			Inclusions: []string{},
			Configs:    configs,
		}, ctx)
	assert.NotEmpty(t, g.Generation.Jobs())
	assert.ErrorIs(t, err, errAssertion)
}

func TestGenerateDirWithoutConfig(t *testing.T) {
	g := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	ctx, _ := ctxTestdata.NewContextMock()
	err := g.Generate(
		DebrickedOptions{
			Paths:      []string{"invalid-dir"},
			Exclusions: []string{},
			Inclusions: []string{},
			Configs:    nil,
		}, ctx)
	assert.Empty(t, g.Generation.Jobs())
	assert.NoError(t, err)
}
