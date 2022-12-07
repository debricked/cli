package scan

import (
	"fmt"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/scan"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

func TestNewScanCmd(t *testing.T) {
	var c client.IDebClient
	c = client.NewDebClient(nil)
	cmd := NewScanCmd(&c)

	viperKeys := viper.AllKeys()
	flags := cmd.Flags()
	flagAssertions := map[string]string{
		RepositoryFlag:    "r",
		CommitFlag:        "c",
		BranchFlag:        "b",
		CommitAuthorFlag:  "a",
		RepositoryUrlFlag: "u",
		IntegrationFlag:   "i",
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

func TestRunE(t *testing.T) {
	var s scan.IScanner
	s = &scannerMock{}
	runE := RunE(&s)
	err := runE(nil, []string{"."})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestRunENoPath(t *testing.T) {
	var s scan.IScanner
	s = &scannerMock{}
	runE := RunE(&s)
	err := runE(nil, []string{})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestRunEError(t *testing.T) {
	runE := RunE(nil)
	err := runE(nil, []string{"."})
	if err == nil {
		t.Error("failed to assert that an error occurred. Error:", err)
	}
	if !strings.Contains(err.Error(), "⨯ scanner was nil") {
		t.Error("failed to assert error message")
	}
}

type scannerMock struct{}

func (*scannerMock) Scan(_ scan.IOptions) error {
	return nil
}
