package license

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/cmd/report/testdata"
	"github.com/debricked/cli/pkg/report"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewLicenseCmd(t *testing.T) {
	var r report.IReporter
	cmd := NewLicenseCmd(r)
	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)

	viperKeys := viper.AllKeys()
	flags := cmd.Flags()
	flagAssertions := map[string]string{
		CommitFlag: "c",
		EmailFlag:  "e",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equal(t, shorthand, flag.Shorthand)

		match := false
		for _, key := range viperKeys {
			if key == name {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that flag was present: "+name)
	}
}

func TestRunEError(t *testing.T) {
	email = "noreply@debricked.com"
	reporterMock := testdata.NewReporterMock()
	reporterMock.SetError(errors.New(""))
	runeE := RunE(reporterMock)

	err := runeE(nil, nil)

	assert.ErrorContains(t, err, "⨯")
}

func TestRunE(t *testing.T) {
	email = "noreply@debricked.com"
	reporterMock := testdata.NewReporterMock()
	runeE := RunE(reporterMock)

	err := runeE(nil, nil)

	assert.NoError(t, err)
}

func TestPreRun(t *testing.T) {
	cmd := NewLicenseCmd(nil)
	cmd.PreRun(cmd, nil)
}
