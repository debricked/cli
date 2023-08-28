package scan

import (
	"testing"

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
		CallGraphUploadTimeoutFlag:   "",
		CallGraphGenerateTimeoutFlag: "",
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

	assert.Error(t, err, scan.FailPipelineErr)
	assert.True(t, cmd.SilenceUsage, "failed to assert that usage was silenced")
	assert.True(t, cmd.SilenceErrors, "failed to assert that errors were silenced")
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
}

func (s *scannerMock) Scan(_ scan.IOptions) error {
	return s.err
}

func (s *scannerMock) setErr(err error) {
	s.err = err
}
