package git

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewMetaObjectWithoutRepositoryName(t *testing.T) {
	metaObj, err := NewMetaObject(".", "", "", "", "", "")

	assert.NotNil(t, metaObj)
	assert.ErrorContains(t, err, "failed to find repository name")
}

func TestNewMetaObjectWithoutCommit(t *testing.T) {
	metaObj, err := NewMetaObject(".", "repository-name", "", "", "", "")

	assert.NotNil(t, metaObj)
	assert.ErrorContains(t, err, "failed to find commit hash")
}

func TestNewMetaObjectWithoutHead(t *testing.T) {
	cwd := testdata.SetUpGitRepository(t, false)
	defer testdata.TearDownGitRepository(cwd, t)

	metaObj, err := NewMetaObject(".", "repository-name", "", "", "", "")

	assert.NotNil(t, metaObj)
	assert.ErrorContains(t, err, "failed to find commit hash")
}

func TestNewMetaObjectWithRepository(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	newMetaObj, err := NewMetaObject(cwd+"/../..", "", "", "", "", "")

	assert.NoError(t, err)
	assert.Equal(t, "debricked/cli", newMetaObj.RepositoryName)
	assert.Equal(t, "https://github.com/debricked/cli", newMetaObj.RepositoryUrl)
	assert.Greater(t, len(newMetaObj.CommitName), 0)
	assert.Greater(t, len(newMetaObj.BranchName), 0)
	assert.Greater(t, len(newMetaObj.DefaultBranchName), 0)
	assert.Greater(t, len(newMetaObj.Author), 0)

}
