package sbom

import (
	"errors"
	"testing"

	"github.com/debricked/cli/internal/cmd/report/testdata"
	"github.com/debricked/cli/internal/report"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewSBOMCmd(t *testing.T) {
	var r report.IReporter
	cmd := NewSBOMCmd(r)
	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)

	viperKeys := viper.AllKeys()
	flags := cmd.Flags()
	flagAssertions := map[string]string{
		CommitFlag:      "c",
		RepositorylFlag: "r",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equalf(t, shorthand, flag.Shorthand, "failed to assert that %s flag shorthand %s was set correctly", name, shorthand)

		match := false
		for _, key := range viperKeys {
			if key == name {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that %s was present", name)
	}
}

func TestRunEError(t *testing.T) {
	reporterMock := testdata.NewReporterMock()
	reporterMock.SetError(errors.New(""))
	runeE := RunE(reporterMock)

	err := runeE(nil, nil)

	assert.ErrorContains(t, err, "⨯")
}

func TestRunE(t *testing.T) {
	reporterMock := testdata.NewReporterMock()
	runeE := RunE(reporterMock)

	err := runeE(nil, nil)

	assert.NoError(t, err)
}

func TestPreRun(t *testing.T) {
	cmd := NewSBOMCmd(nil)
	cmd.PreRun(cmd, nil)
}
