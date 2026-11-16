package find

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/debricked/cli/internal/cmd/cmderror"
	sanitizecmd "github.com/debricked/cli/internal/cmd/sanitize"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/file/testdata"
	"github.com/debricked/cli/internal/sanitize"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewFindCmd(t *testing.T) {
	var f file.IFinder
	cmd := NewFindCmd(f, nil)

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
	runE := RunE(f, nil)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunENoPath(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	runE := RunE(f, nil)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunENoFiles(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{})
	f.SetGetGroupsReturnMock(groups, nil)
	exclusions = []string{}
	runE := RunE(f, nil)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunEError(t *testing.T) {
	f := testdata.NewFinderMock()
	errorAssertion := errors.New("finder-error")
	f.SetGetGroupsReturnMock(file.Groups{}, errorAssertion)
	runE := RunE(f, nil)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "finder-error", "error doesn't match expected")
}

func TestRunEWithInvalidStrictFlag(t *testing.T) {
	viper.Set(StrictFlag, 123)

	f := testdata.NewFinderMock()
	runE := RunE(f, nil)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "'strict' supports values within range 0-2", "error doesn't match expected")
}

func TestRunEWithBothStrictAndLockOnlyFlagsSet(t *testing.T) {
	viper.Set(StrictFlag, file.StrictLockAndPairs)
	viper.Set(LockfileOnlyFlag, true)

	defer viper.Set(StrictFlag, file.StrictAll)
	defer viper.Set(LockfileOnlyFlag, false)

	f := testdata.NewFinderMock()
	runE := RunE(f, nil)
	err := runE(nil, []string{"."})

	assert.EqualError(t, err, "'lockfile' and 'strict' flags are mutually exclusive", "error doesn't match expected")
}

func TestRunEJson(t *testing.T) {
	f := testdata.NewFinderMock()
	groups := file.Groups{}
	groups.Add(file.Group{
		ManifestFile:   "manifest-file",
		CompiledFormat: nil,
		LockFiles: []string{
			"lock-file",
		},
	})
	groupsJson, _ := json.Marshal(groups.ToSlice())
	f.SetGetGroupsReturnMock(groups, nil)
	viper.Set(JsonFlag, true)

	runE := RunE(f, nil)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runE(nil, []string{"."})
	assert.NoError(t, err)

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.NotEmpty(t, output)
	assert.JSONEq(t, string(groupsJson), string(output))
}

func TestPreRun(t *testing.T) {
	cmd := NewFindCmd(nil, nil)
	cmd.PreRun(cmd, nil)
}

type sanitizerMock struct {
	result  sanitize.Result
	options sanitize.Options
}

func (m *sanitizerMock) Sanitize(options sanitize.Options) (sanitize.Result, error) {
	m.options = options

	return m.result, nil
}

func TestRunESanitize(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set(SanitizeFlag, true)
	viper.Set(LockfileOnlyFlag, true)
	viper.Set(ExclusionFlag, []string{"**/node_modules/**"})
	finder := testdata.NewFinderMock()
	finder.SetGetGroupsReturnMock(file.Groups{}, errors.New("the finder must not be called directly"))
	sanitizer := &sanitizerMock{result: sanitize.Result{Groups: []file.Group{
		*file.NewGroup("", nil, []string{".debricked/sanitized/yarn.lock"}),
	}}}
	cmd := NewFindCmd(finder, sanitizer)
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)

	err := RunE(finder, sanitizer)(cmd, []string{"app"})

	assert.NoError(t, err)
	assert.Equal(t, sanitize.Options{
		RootPath:     "app",
		Exclusions:   []string{"**/node_modules/**"},
		LockFileOnly: true,
		Strictness:   file.StrictAll,
	}, sanitizer.options)
	assert.Equal(t, ".debricked/sanitized/yarn.lock\n", out.String())
	assert.Contains(t, errOut.String(), "Sanitized: 0, skipped: 0, failed: 0")
}

func TestRunESanitizeReturnsTheExitCodeOfSanitize(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set(SanitizeFlag, true)
	sanitizer := &sanitizerMock{result: sanitize.Result{
		Failed: []sanitize.FileIssue{{Path: "yarn.lock", Reason: "content is not valid for this file type"}},
	}}
	cmd := NewFindCmd(nil, sanitizer)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := RunE(nil, sanitizer)(cmd, nil)

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, sanitizecmd.FailureExitCode, cmdErr.Code)
}
