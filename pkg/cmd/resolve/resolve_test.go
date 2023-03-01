package resolve

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/file/testdata"
	"github.com/debricked/cli/pkg/resolution"
	resolveTestdata "github.com/debricked/cli/pkg/resolution/testdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewResolveCmd(t *testing.T) {
	var resolver resolution.IResolver
	cmd := NewResolveCmd(resolver)

	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)

	flags := cmd.Flags()
	flagAssertions := map[string]string{}
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
	f := testdata.NewFinderMock()
	r := &resolveTestdata.ResolverMock{}
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	runE := RunE(r)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunENoPath(t *testing.T) {
	f := testdata.NewFinderMock()
	r := &resolveTestdata.ResolverMock{}
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	runE := RunE(r)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunENoFiles(t *testing.T) {
	f := testdata.NewFinderMock()
	r := &resolveTestdata.ResolverMock{}
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	exclusions = []string{}
	runE := RunE(r)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunEError(t *testing.T) {
	r := &resolveTestdata.ResolverMock{}
	errorAssertion := errors.New("finder-error")
	r.Err = errorAssertion
	runE := RunE(r)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "finder-error", "error doesn't match expected")
}
