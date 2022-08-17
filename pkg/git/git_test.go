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
		t.Fatal("failed to get repository. Error:", err)
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

func TestFindRepositoryNameWithoutMetaData(t *testing.T) {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	store := memory.NewStorage()
	r, err := git.Init(store, fs)
	if err != nil {
		t.Fatal("failed to get repository. Error:", err)
	}
	name, err := FindRepositoryName(r, "test/repository")
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if name != "repository" {
		t.Error("failed to find repository name")
	}
}
