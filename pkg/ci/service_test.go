package ci

import (
	"github.com/debricked/cli/pkg/ci/azure"
	"github.com/debricked/cli/pkg/ci/circleci"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	s := NewService([]ICi{})
	if len(s.cis) > 0 {
		t.Error("failed to assert that CiService lacked CIs")
	}
	s = NewService(nil)
	if len(s.cis) != 8 {
		t.Error("failed to assert number of CIs")
	}
}

func TestFindNotSupported(t *testing.T) {
	s := NewService([]ICi{azure.Ci{}, circleci.Ci{}})
	_, err := s.Find()
	if err != ErrNotSupported {
		t.Error("failed to assert that error was ErrNotSupported")
	}
}

func TestFind(t *testing.T) {
	_ = os.Setenv(azure.EnvKey, "value")
	defer os.Unsetenv(azure.EnvKey)

	s := NewService([]ICi{azure.Ci{}, circleci.Ci{}})
	env, err := s.Find()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}

	if len(env.Integration) == 0 {
		t.Error("failed to assert that was was no value in env")
	}
}
