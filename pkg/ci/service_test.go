package ci

import (
	"os"
	"testing"

	"github.com/debricked/cli/pkg/ci/azure"
	"github.com/debricked/cli/pkg/ci/circleci"
	"github.com/debricked/cli/pkg/ci/gitlab"
	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	s := NewService([]ICi{})
	assert.Empty(t, s.cis)

	s = NewService(nil)
	assert.Len(t, s.cis, 8)

	s.cis = []ICi{gitlab.Ci{}}
	assert.Len(t, s.cis, 1)

	_, ok := s.cis[0].(gitlab.Ci)
	assert.True(t, ok, "failed to assert that the CI was gitlab.Ci")
}

func TestFindNotSupported(t *testing.T) {
	s := NewService([]ICi{azure.Ci{}, circleci.Ci{}})
	_, err := s.Find()
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestFind(t *testing.T) {
	_ = os.Setenv(azure.EnvKey, "value")
	defer testdata.UnsetEnvVar(t, azure.EnvKey)

	s := NewService([]ICi{azure.Ci{}, circleci.Ci{}})
	env, err := s.Find()
	assert.NoError(t, err)
	assert.Greater(t, len(env.Integration), 0)
}
