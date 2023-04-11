package strategy

import (
	"testing"

	"github.com/debricked/cli/pkg/resolution/file"
	"github.com/debricked/cli/pkg/resolution/pm/gomod"
	"github.com/debricked/cli/pkg/resolution/pm/gradle"
	"github.com/debricked/cli/pkg/resolution/pm/maven"
	"github.com/debricked/cli/pkg/resolution/pm/pip"
	"github.com/debricked/cli/pkg/resolution/pm/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategyFactory(t *testing.T) {
	f := NewStrategyFactory()
	assert.NotNil(t, f)
}

func TestMakeErr(t *testing.T) {
	f := NewStrategyFactory()
	batch := file.NewBatch(testdata.PmMock{N: "test"})
	s, err := f.Make(batch, nil)
	assert.Nil(t, s)
	assert.ErrorContains(t, err, "failed to make strategy from test")
}

func TestMake(t *testing.T) {
	cases := map[string]IStrategy{
		maven.Name:  maven.NewStrategy(nil),
		gradle.Name: gradle.NewStrategy(nil, nil),
		gomod.Name:  gomod.NewStrategy(nil),
		pip.Name:    pip.NewStrategy(nil),
	}
	f := NewStrategyFactory()
	var batch file.IBatch
	for name, strategy := range cases {
		batch = file.NewBatch(testdata.PmMock{N: name})
		t.Run(name, func(t *testing.T) {
			s, err := f.Make(batch, nil)
			assert.NoError(t, err)
			assert.Equal(t, strategy, s)
		})
	}
}
