package docker

import (
	"errors"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/scan"
	"strings"
	"testing"
)

func TestNewDockerCmd(t *testing.T) {
	var c client.IDebClient
	c = client.NewDebClient(nil)
	cmd := NewDockerCmd(&c)
	if cmd == nil {
		t.Error("Failed to assert that docker command was created")
	}
}

func TestRunE(t *testing.T) {
	var s scan.IScanner
	s = &scannerMock{}
	runE := RunE(s)
	err := runE(nil, []string{"debricked/cli"})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestRunENoPath(t *testing.T) {
	var s scan.IScanner
	s = &scannerMock{}
	runE := RunE(s)
	_ = runE(nil, []string{})
}

func TestRunEError(t *testing.T) {
	var s scan.IScanner
	s = &scannerMock{err: true}
	runE := RunE(s)
	err := runE(nil, []string{"debricked/cli"})
	if err == nil {
		t.Error("failed to assert that an error occurred. Error:", err)
	}
	if !strings.Contains(err.Error(), "error") {
		t.Error("failed to assert error message")
	}
}

func TestValidateArgs(t *testing.T) {
	err := ValidateArgs(nil, []string{})
	if !strings.EqualFold("failed to validate argument. Please use one argument", err.Error()) {
		t.Error("failed to validate error message")
	}
	err = ValidateArgs(nil, []string{"debricked/cli"})
	if err != nil {
		t.Error("failed to asser that no error occurred")
	}
	err = ValidateArgs(nil, []string{"debricked/cli", "debricked/scan"})
	if !strings.EqualFold("failed to validate argument. Please use one argument", err.Error()) {
		t.Error("failed to validate error message")
	}
}

type scannerMock struct {
	err bool
}

func (s *scannerMock) Scan(_ scan.IOptions) error {
	if s.err {
		return errors.New("error")
	}
	return nil
}
