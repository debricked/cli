package git

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"os"
	"strings"
	"testing"
)

func TestNewMetaObjectWithoutRepositoryName(t *testing.T) {
	metaObj, err := NewMetaObject(".", "", "", "", "", "")
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if metaObj == nil {
		t.Error("failed to assert that gitMetaObject was not nil")
	}
	if !strings.Contains(err.Error(), "failed to find repository name") {
		t.Error("failed to assert that repository name was missing")
	}
}

func TestNewMetaObjectWithoutCommit(t *testing.T) {
	metaObj, err := NewMetaObject(".", "repository-name", "", "", "", "")
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if metaObj == nil {
		t.Error("failed to assert that gitMetaObject was not nil")
	}
	if !strings.Contains(err.Error(), "failed to find commit hash") {
		t.Error("failed to assert that commit hash was missing")
	}
}

func TestNewMetaObjectWithoutHead(t *testing.T) {
	cwd, err := testdata.SetUpGitRepository(false)
	if err != nil {
		t.Fatal("failed to initialize repository", err)
	}
	defer testdata.TearDownGitRepository(cwd)

	metaObj, err := NewMetaObject(".", "repository-name", "", "", "", "")
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if metaObj == nil {
		t.Error("failed to assert that gitMetaObject was not nil")
	}
	if !strings.Contains(err.Error(), "failed to find commit hash") {
		t.Error("failed to assert that commit hash was missing")
	}
}

func TestNewMetaObjectWithRepository(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	newMetaObj, err := NewMetaObject(cwd+"/../..", "", "", "", "", "")
	if err != nil {
		t.Error(err)
	}
	if newMetaObj.RepositoryName != "debricked/cli" {
		t.Error("failed to find correct repository name:", newMetaObj.RepositoryName)
	}
	if newMetaObj.RepositoryUrl != "https://github.com/debricked/cli" {
		t.Error("failed to find correct repository url:", newMetaObj.RepositoryUrl)
	}
	if len(newMetaObj.CommitName) == 0 {
		t.Error("failed to find correct commit", newMetaObj.CommitName)
	}
	if len(newMetaObj.BranchName) == 0 || len(newMetaObj.DefaultBranchName) == 0 {
		t.Error("failed to find correct branch", newMetaObj.BranchName)
	}
	if len(newMetaObj.Author) == 0 {
		t.Error("failed to find correct commit author", newMetaObj.Author)
	}
}
