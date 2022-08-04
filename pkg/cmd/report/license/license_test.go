package license

import (
	"debricked/pkg/client"
	"fmt"
	"strings"
	"testing"
)

const validCommit = "b3b1ff886344d876d13ab916bcfdba41c4e7a8bb"

func TestNewLicenseCmd(t *testing.T) {
	var c client.Client = client.NewDebClient(nil)
	cmd := NewLicenseCmd(&c)
	commands := cmd.Commands()
	nbrOfCommands := 0
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}

	flags := cmd.Flags()
	flagAssertions := map[string]string{
		"commit": "c",
		"email":  "e",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		if flag == nil {
			t.Error(fmt.Sprintf("failed to assert that %s flag was set", name))
		}
		if flag.Shorthand != shorthand {
			t.Error(fmt.Sprintf("failed to assert that %s flag shorthand %s was set correctly", name, shorthand))
		}
	}
}

func TestRunUnAuthorized(t *testing.T) {
	email = "noreply@debricked.com"
	commitHash = validCommit
	accessToken := "invalid"
	debClient = client.NewDebClient(&accessToken)
	err := run(nil, nil)
	if err == nil {
		t.Fatal("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "⨯ Unauthorized. Specify access token") {
		t.Error("failed to assert error message")
	}
}

func TestRun(t *testing.T) {
	email = "noreply@debricked.com"
	commitHash = validCommit
	debClient = client.NewDebClient(nil)
	err := run(nil, nil)
	if err != nil {
		t.Fatal("failed to assert that no error occurred")
	}
}

func TestReportInvalidCommitHash(t *testing.T) {
	email = "noreply@debricked.com"
	commitHash = "invalid"
	debClient = client.NewDebClient(nil)
	err := report()
	if err == nil {
		t.Fatal("failed to assert that error occurred")
	}
	if !strings.Contains(err.Error(), "No commit was found with the name invalid") {
		t.Error("failed to assert error message")
	}
}

func TestGetCommitId(t *testing.T) {
	debClient = client.NewDebClient(nil)
	id, err := getCommitId(validCommit)
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
	if id < 1 {
		t.Error("failed to assert that the commit ID was a positive integer")
	}
}
