package callgraph

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/callgraph/config"
	strategyTestdata "github.com/debricked/cli/pkg/callgraph/strategy/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	workers   = 10
	goModFile = "go.mod"
)

func TestNewGenerator(t *testing.T) {
	r := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)
	assert.NotNil(t, r)
}

func TestGenerate(t *testing.T) {
	r := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		NewScheduler(workers),
	)

	var status chan bool
	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	res, err := r.Generate([]string{"../../go.mod"}, nil, configs, status)
	assert.NotEmpty(t, res.Jobs())
	assert.NoError(t, err)
}

func TestGenerateInvokeError(t *testing.T) {
	r := NewGenerator(
		strategyTestdata.NewStrategyFactoryErrorMock(),
		NewScheduler(workers),
	)

	var status chan bool
	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	_, err := r.Generate([]string{"../../go.mod"}, nil, configs, status)
	assert.NotNil(t, err)
}

func TestGenerateScheduleError(t *testing.T) {
	errAssertion := errors.New("error")
	r := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{Err: errAssertion},
	)

	var status chan bool
	configs := []config.IConfig{
		config.NewConfig("java", []string{}, map[string]string{"pm": "maven"}),
	}
	res, err := r.Generate([]string{"../../go.mod"}, nil, configs, status)
	assert.NotEmpty(t, res.Jobs())
	assert.ErrorIs(t, err, errAssertion)
}

func TestGenerateDirWithoutConfig(t *testing.T) {
	r := NewGenerator(
		strategyTestdata.NewStrategyFactoryMock(),
		SchedulerMock{},
	)

	var status chan bool
	res, err := r.Generate([]string{"invalid-dir"}, nil, nil, status)
	assert.Empty(t, res.Jobs())
	assert.NoError(t, err)
}
