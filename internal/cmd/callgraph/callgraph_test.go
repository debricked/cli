package callgraph

import (
	"errors"
	"testing"

	"github.com/debricked/cli/internal/callgraph"
	callgraphTestdata "github.com/debricked/cli/internal/callgraph/testdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewCallgraphCmd(t *testing.T) {
	var callgraphGenerator callgraph.IGenerator
	cmd := NewCallgraphCmd(callgraphGenerator)

	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)

	flags := cmd.Flags()
	flagAssertions := map[string]string{
		ExclusionFlag:       "e",
		InclusionFlag:       "",
		NoBuildFlag:         "",
		GenerateTimeoutFlag: "",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equal(t, shorthand, flag.Shorthand)
	}

	var flagKeys = []string{
		ExclusionFlag,
	}
	viperKeys := viper.AllKeys()
	for _, flagKey := range flagKeys {
		match := false
		for _, key := range viperKeys {
			if key == flagKey {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that flag was present: "+flagKey)
	}

}

func TestRunE(t *testing.T) {
	g := &callgraphTestdata.GeneratorMock{}
	runE := RunE(g)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunENoPath(t *testing.T) {
	g := &callgraphTestdata.GeneratorMock{}
	runE := RunE(g)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunEError(t *testing.T) {
	g := &callgraphTestdata.GeneratorMock{}
	errorAssertion := errors.New("finder-error")
	g.Err = errorAssertion
	runE := RunE(g)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "finder-error", "error doesn't match expected")

	languages = "python2"

	g2 := &callgraphTestdata.GeneratorMock{}
	runE2 := RunE(g2)
	err2 := runE2(nil, []string{"."})
	assert.Error(t, err2)

}

func TestParseAndValidateLanguages(t *testing.T) {
	languages := "java,golang"
	parsedLanguages, err := parseAndValidateLanguages(languages)

	assert.Nil(t, err)
	assert.Equal(t, []string{"java", "golang"}, parsedLanguages)

	languages = ""
	parsedLanguages, err = parseAndValidateLanguages(languages)

	assert.Nil(t, err)
	assert.Equal(t, []string{"java", "golang"}, parsedLanguages)

	languages = "java,golang,python2"
	_, err = parseAndValidateLanguages(languages)
	assert.Error(t, err)
}
