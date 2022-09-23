package license

import (
	"debricked/pkg/cmd/report/testdata"
	"debricked/pkg/report"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

func TestNewLicenseCmd(t *testing.T) {
	var r report.IReporter
	cmd := NewLicenseCmd(r)
	commands := cmd.Commands()
	nbrOfCommands := 0
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}

	viperKeys := viper.AllKeys()
	flags := cmd.Flags()
	flagAssertions := map[string]string{
		CommitFlag: "c",
		EmailFlag:  "e",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		if flag == nil {
			t.Error(fmt.Sprintf("failed to assert that %s flag was set", name))
		}
		if flag.Shorthand != shorthand {
			t.Error(fmt.Sprintf("failed to assert that %s flag shorthand %s was set correctly", name, shorthand))
		}

		match := false
		for _, key := range viperKeys {
			if key == name {
				match = true
			}
		}
		if !match {
			t.Error("failed to assert that flag was present: " + name)
		}
	}
}

func TestRunEError(t *testing.T) {
	email = "noreply@debricked.com"
	reporterMock := testdata.NewReporterMock()
	reporterMock.SetError(errors.New(""))
	runeE := RunE(reporterMock)
	err := runeE(nil, nil)
	if !strings.Contains(err.Error(), "⨯") {
		t.Error("failed to assert error message")
	}
}

func TestRunE(t *testing.T) {
	email = "noreply@debricked.com"
	reporterMock := testdata.NewReporterMock()
	runeE := RunE(reporterMock)
	err := runeE(nil, nil)
	if err != nil {
		t.Fatal("failed to assert that no error occurred")
	}
}
