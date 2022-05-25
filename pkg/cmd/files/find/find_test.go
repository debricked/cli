package find

import (
	"debricked/pkg/client"
	"debricked/pkg/file"
	"fmt"
	"strings"
	"testing"
)

func TestNewFindCmd(t *testing.T) {
	cmd := NewFindCmd(client.NewDebClient(nil))

	commands := cmd.Commands()
	nbrOfCommands := 0
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}

	flags := cmd.Flags()
	flagAssertions := map[string]string{}
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

func TestRun(t *testing.T) {
	debClient = client.NewDebClient(nil)
	finder, _ = file.NewFinder(debClient)
	err := run(nil, []string{"."})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestFind(t *testing.T) {
	finder, _ = file.NewFinder(client.NewDebClient(nil))
	err := find("../../scan", []string{})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestValidateArgs(t *testing.T) {
	err := validateArgs(nil, []string{"."})
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
}

func TestValidateArgsInvalidArgs(t *testing.T) {
	err := validateArgs(nil, []string{})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "requires path") {
		t.Error("failed to assert error message")
	}

	err = validateArgs(nil, []string{"invalid-path"})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "invalid path specified") {
		t.Error("failed to assert error message")
	}
}
