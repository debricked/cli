package ci

import (
	"debricked/pkg/ci/circleci"
	"debricked/pkg/ci/gitlab"
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
	s := NewService([]ICi{gitlab.Ci{}, circleci.Ci{}})
	_, err := s.Find()
	if err != ErrNotSupported {
		t.Error("failed to assert that error was ErrNotSupported")
	}
}

func TestFind(t *testing.T) {
	_ = os.Setenv(circleci.EnvKey, "value")
	defer os.Unsetenv(circleci.EnvKey)

	s := NewService([]ICi{gitlab.Ci{}, circleci.Ci{}})
	env, err := s.Find()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}

	if len(env.Repository) != 0 {
		t.Error("failed to assert that was was no value in env")
	}
}
