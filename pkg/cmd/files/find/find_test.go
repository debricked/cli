package find

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/file/testdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewFindCmd(t *testing.T) {
	var f file.IFinder
	cmd := NewFindCmd(f)

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
		JsonFlag,
		LockfileOnlyFlag,
		StrictFlag,
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
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	runE := RunE(f)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunENoPath(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	runE := RunE(f)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunENoFiles(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	exclusions = []string{}
	runE := RunE(f)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunEError(t *testing.T) {
	f := testdata.NewFinderMock()
	errorAssertion := errors.New("finder-error")
	f.SetGetGroupsReturnMock(file.Groups{}, errorAssertion)
	runE := RunE(f)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "finder-error", "error doesn't match expected")
}

func TestRunEWithInvalidStrictFlag(t *testing.T) {
	viper.Set(StrictFlag, 123)

	f := testdata.NewFinderMock()
	runE := RunE(f)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "'strict' supports values within range 0-2", "error doesn't match expected")
}

func TestRunEWithBothStrictAndLockOnlyFlagsSet(t *testing.T) {
	viper.Set(StrictFlag, file.StrictLockAndPairs)
	viper.Set(LockfileOnlyFlag, true)

	f := testdata.NewFinderMock()
	runE := RunE(f)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "'lockfile' and 'strict' flags are mutually exclusive", "error doesn't match expected")
}
