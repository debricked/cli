package scan

import (
	"errors"
	"testing"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/resolution"
	"github.com/debricked/cli/internal/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewScanCmd(t *testing.T) {
	cmd := NewScanCmd(&scannerMock{})

	flagAssertions := map[string]string{
		RepositoryFlag:               "r",
		CommitFlag:                   "c",
		BranchFlag:                   "b",
		CommitAuthorFlag:             "a",
		RepositoryUrlFlag:            "u",
		IntegrationFlag:              "i",
		ExclusionFlag:                "e",
		PassOnTimeOut:                "p",
		NoResolveFlag:                "",
		CallGraphFlag:                "",
		JavaCallgraphEngineFlag:      "",
		CallGraphUploadTimeoutFlag:   "",
		CallGraphGenerateTimeoutFlag: "",
		ResolutionStrictnessFlag:     "",
	}
	flags := cmd.Flags()
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equalf(t, shorthand, flag.Shorthand, "failed to assert that %s flag shorthand %s was set correctly", name, shorthand)
	}

	var flagKeys = []string{
		RepositoryFlag,
		CommitFlag,
		BranchFlag,
		CommitAuthorFlag,
		RepositoryUrlFlag,
		IntegrationFlag,
		TagCommitAsReleaseFlag,
	}
	viperKeys := viper.AllKeys()
	for _, flagKey := range flagKeys {
		match := false
		for _, key := range viperKeys {
			if key == flagKey {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that %s was present", flagKey)
	}
}

func TestRunE(t *testing.T) {
	var s scan.IScanner = &scannerMock{}
	runE := RunE(&s)

	err := runE(nil, []string{"."})

	assert.NoError(t, err)
}

func TestRunENoPath(t *testing.T) {
	var s scan.IScanner = &scannerMock{}
	runE := RunE(&s)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunEFailPipelineErr(t *testing.T) {
	var s scan.IScanner
	mock := &scannerMock{}
	mock.setErr(scan.FailPipelineErr)
	s = mock
	runE := RunE(&s)
	cmd := &cobra.Command{}

	err := runE(cmd, nil)

	assert.ErrorIs(t, err, scan.FailPipelineErr)
	assert.True(t, cmd.SilenceUsage, "failed to assert that usage was silenced")
	assert.True(t, cmd.SilenceErrors, "failed to assert that errors were silenced")
}

func TestRunELongQueueErr(t *testing.T) {
	var s scan.IScanner
	mock := &scannerMock{}
	mock.setErr(scan.LongQueueErr)
	s = mock
	runE := RunE(&s)
	cmd := &cobra.Command{}

	err := runE(cmd, nil)

	assert.ErrorIs(t, err, scan.LongQueueErr)
	assert.True(t, cmd.SilenceUsage, "failed to assert that usage was silenced")
	assert.True(t, cmd.SilenceErrors, "failed to assert that errors were silenced")
}

func TestRunECommandError(t *testing.T) {
	var s scan.IScanner
	mock := &scannerMock{}
	cmdErr := cmderror.CommandError{Code: 3, Err: errors.New("partial resolution failure")}
	mock.setErr(cmdErr)
	s = mock
	runE := RunE(&s)
	cmd := &cobra.Command{}

	err := runE(cmd, nil)

	var gotCmdErr cmderror.CommandError
	assert.True(t, errors.As(err, &gotCmdErr), "expected CommandError to be preserved")
	assert.Equal(t, 3, gotCmdErr.Code, "expected exit code 3 to be preserved")
	assert.True(t, cmd.SilenceUsage, "failed to assert that usage was silenced")
	assert.True(t, cmd.SilenceErrors, "failed to assert that errors were silenced")
}

func TestRunEResolutionStrictness(t *testing.T) {
	cases := []struct {
		name     string
		flag     interface{}
		expected resolution.StrictnessLevel
	}{
		{name: "default", flag: nil, expected: resolution.FailIfAllFail},
		{name: "no fail", flag: 0, expected: resolution.NoFail},
		{name: "fail if all fail", flag: 1, expected: resolution.FailIfAllFail},
		{name: "fail if any fail", flag: 2, expected: resolution.FailIfAnyFail},
		{name: "fail or warn", flag: 3, expected: resolution.FailOrWarn},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			viper.Reset()
			if c.flag != nil {
				viper.Set(ResolutionStrictnessFlag, c.flag)
			} else {
				// Mirror the flag default that PreRun would have bound.
				viper.SetDefault(ResolutionStrictnessFlag, int(resolution.FailIfAllFail))
			}
			defer viper.Reset()

			var s scan.IScanner
			mock := &scannerMock{}
			s = mock
			runE := RunE(&s)

			err := runE(&cobra.Command{}, nil)

			assert.NoError(t, err)
			options, ok := mock.options.(scan.DebrickedOptions)
			assert.True(t, ok, "failed to assert that scan options were passed")
			assert.Equal(t, c.expected, options.ResolutionStrictness)
		})
	}
}

// "debricked files find" binds DEBRICKED_STRICT to the global viper key
// "strict". The scan flag must not share that key, or setting the env var for
// one command would silently change resolution behaviour in the other.
func TestRunEStrictEnvDoesNotAffectResolutionStrictness(t *testing.T) {
	viper.Reset()
	viper.SetEnvPrefix("DEBRICKED")
	viper.AutomaticEnv()
	viper.MustBindEnv("strict")
	viper.SetDefault(ResolutionStrictnessFlag, int(resolution.FailIfAllFail))
	t.Setenv("DEBRICKED_STRICT", "3")
	defer viper.Reset()

	var s scan.IScanner
	mock := &scannerMock{}
	s = mock
	runE := RunE(&s)

	err := runE(&cobra.Command{}, nil)

	assert.NoError(t, err)
	options, ok := mock.options.(scan.DebrickedOptions)
	assert.True(t, ok, "failed to assert that scan options were passed")
	assert.Equal(t, resolution.FailIfAllFail, options.ResolutionStrictness)
}

func TestRunEInvalidResolutionStrictness(t *testing.T) {
	viper.Reset()
	viper.Set(ResolutionStrictnessFlag, 4)
	defer viper.Reset()

	var s scan.IScanner
	mock := &scannerMock{}
	s = mock
	runE := RunE(&s)

	err := runE(&cobra.Command{}, nil)

	assert.ErrorContains(t, err, "invalid strictness level: 4")
	assert.Nil(t, mock.options, "failed to assert that the scan was not started")
}

func TestRunEError(t *testing.T) {
	runE := RunE(nil)
	err := runE(nil, []string{"."})

	assert.ErrorContains(t, err, "⨯ scanner was nil")
}

func TestPreRun(t *testing.T) {
	cmd := NewScanCmd(nil)
	cmd.PreRun(cmd, nil)
}

type scannerMock struct {
	err error
	// options records the options of the most recent Scan call.
	options scan.IOptions
}

func (s *scannerMock) Scan(o scan.IOptions) error {
	s.options = o

	return s.err
}

func (s *scannerMock) setErr(err error) {
	s.err = err
}
