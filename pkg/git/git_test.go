package git

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
	"strings"
	"testing"
)

var repository *git.Repository

func setUp(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	repository, err = git.PlainOpen(cwd + "/../..")
	if err != nil {
		t.Error(err)
	}
}

func TestNewMetaObjectWithoutRepositoryName(t *testing.T) {
	metaObj, err := NewMetaObject(".", "", "", "", "", "")
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if metaObj == nil {
		t.Error("failed to assert that gitMetaObject was not nil")
	}
	if !strings.Contains(err.Error(), "failed to find repository name. Please use --repository flag") {
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
	if !strings.Contains(err.Error(), "failed to find commit hash. Please use --commit flag") {
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

func TestFindBranchName(t *testing.T) {
	setUp(t)
	head, _ := repository.Head()
	branch, err := FindBranchName(repository, head.Hash().String())
	if err != nil {
		t.Error(err)
	}
	if len(branch) == 0 {
		t.Error("failed to find correct branch", branch)
	}
}

func TestFindRepositoryName(t *testing.T) {
	setUp(t)
	name, err := FindRepositoryName(repository, "")
	if err != nil {
		t.Error(err)
	}
	if name != "debricked/cli" {
		t.Error("failed to find correct repository name:", name)
	}
}

func TestFindRepositoryUrl(t *testing.T) {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	store := memory.NewStorage()
	r, err := git.Init(store, fs)
	remoteConfig := &config.RemoteConfig{
		Name:  "debricked/cli",
		URLs:  []string{"git@github.com:debricked/cli.git"},
		Fetch: nil,
	}
	_, err = r.CreateRemote(remoteConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	url, err := FindRepositoryUrl(r)
	if err != nil {
		t.Error(err.Error())
	}
	if url != "https://github.com/debricked/cli" {
		t.Error("failed to find correct repository url:", url)
	}
}

func TestParseGitRemoteUrl(t *testing.T) {
	remoteUrls := []string{
		"git@github.com:debricked/cli.git",
		"https://github.com/debricked/cli.git",
		"https://some.git.host/debricked/cli.git",
		"https://gitlab.com/debricked/cli.git",
		"git@gitlab.com:debricked/cli.git",
		"https://github.com/debricked/cli",
	}
	for _, remoteUrl := range remoteUrls {
		name, err := ParseGitRemoteUrl(remoteUrl)
		if err != nil {
			t.Error(err)
		}
		if name != "debricked/cli" {
			t.Error("failed to find correct git repository name:", name)
		}
	}
	badUrl := "ftp://github.com/debricked/cli"
	name, err := ParseGitRemoteUrl(badUrl)
	if err == nil {
		t.Error("failed to assert that bad URL generated error")
	}
	if name != badUrl {
		t.Error("failed to assert that parsed git remote URL equals inputted URL")
	}
}

func TestFindRemoteUrl(t *testing.T) {
	setUp(t)
	url, err := FindRemoteUrl(repository)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(url, "debricked/cli") || !strings.Contains(url, "github.com") {
		t.Error("failed to find correct git remote url:", url)
	}
}

func TestFindRepositoryUrlWithFailure(t *testing.T) {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	store := memory.NewStorage()
	r, err := git.Init(store, fs)
	if err != nil {
		t.Fatal("failed to get repository. Error: ", err.Error())
	}
	url, err := FindRepositoryUrl(r)
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if url != "" {
		t.Error("failed to assert that url was ")
	}
	if !strings.Contains(err.Error(), "failed to find repository remote URL") {
		t.Error("failed to assert error message")
	}
}

func TestFindBranchNameWithFailure(t *testing.T) {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	store := memory.NewStorage()
	r, err := git.Init(store, fs)
	if err != nil {
		t.Fatal("failed to get repository. Error: ", err.Error())
	}
	branch, err := FindBranchName(r, "commit-hash")
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if branch != "" {
		t.Error("failed to find correct branch")
	}
	if !strings.Contains(err.Error(), "failed to find branch") {
		t.Error("failed to assert error message")
	}
}

func TestFindRepositoryNameWithoutMetaData(t *testing.T) {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	store := memory.NewStorage()
	r, err := git.Init(store, fs)
	if err != nil {
		t.Fatal("failed to get repository. Error: ", err.Error())
	}
	name, err := FindRepositoryName(r, "test/repository")
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if name != "repository" {
		t.Error("failed to find repository name")
	}
}
