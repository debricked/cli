package strategy

import (
	"testing"

	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder/javafinder"
	"github.com/debricked/cli/internal/callgraph/language/java"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategyFactory(t *testing.T) {
	f := NewStrategyFactory()
	assert.NotNil(t, f)
}

func TestMakeErr(t *testing.T) {
	f := NewStrategyFactory()
	conf := config.NewConfig("test", nil, nil, true, "")
	s, err := f.Make(conf, nil, nil, nil, nil)
	assert.Nil(t, s)
	assert.ErrorContains(t, err, "failed to make strategy from test")
}

func TestMake(t *testing.T) {
	conf := config.NewConfig(java.Name, nil, nil, true, "")
	cases := map[string]IStrategy{
		java.Name: java.NewStrategy(conf, []string{}, []string{}, []string{}, javafinder.JavaFinder{}, nil),
	}
	f := NewStrategyFactory()
	for name, strategy := range cases {
		t.Run(name, func(t *testing.T) {
			s, err := f.Make(conf, []string{}, []string{}, []string{}, nil)
			assert.NoError(t, err)
			assert.Equal(t, strategy, s)
		})
	}
}
