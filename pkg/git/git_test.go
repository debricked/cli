package git

import (
	"os"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	debrickedCli = "debricked/cli"
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
	if name != debrickedCli {
		t.Error("failed to find correct repository name:", name)
	}
}

func TestFindRepositoryUrl(t *testing.T) {
	remoteUrls := []string{
		"git@github.com:debricked/cli.git",
		"git@github.com:debricked/cli",
		"ssh://git@github.com/debricked/cli",
		"ssh://git@github.com/debricked/cli.git",
		"https://github.com/debricked/cli.git",
		"https://github.com/debricked/cli",
	}

	for _, remoteUrl := range remoteUrls {
		repoMock := mockRepository(true, t)
		remoteConfig := &config.RemoteConfig{
			Name:  debrickedCli,
			URLs:  []string{remoteUrl},
			Fetch: nil,
		}
		_, err := repoMock.CreateRemote(remoteConfig)
		if err != nil {
			t.Fatal(err.Error())
		}
		url, err := FindRepositoryUrl(repoMock)
		if err != nil {
			t.Error(err.Error())
		}
		if url != "https://github.com/debricked/cli" {
			t.Error("failed to find correct repository url from:", remoteUrl, "got:", url)
		}
	}
}

func TestParseGitRemoteUrl(t *testing.T) {
	remoteUrls := []string{
		"git@github.com:debricked/cli.git",
		"git@github.com:debricked/cli",
		"https://github.com/debricked/cli.git",
		"https://some.git.host/debricked/cli.git",
		"https://gitlab.com/debricked/cli.git",
		"git@gitlab.com:debricked/cli.git",
		"git@some.git.host:debricked/cli.git",
		"https://github.com/debricked/cli",
		"ssh://git@github.com/debricked/cli.git",
		"ssh://git@github.com/debricked/cli",
		"ssh://git@some.git.host/debricked/cli.git",
		"ssh://git@some.git.host:1337/debricked/cli.git",
	}
	for _, remoteUrl := range remoteUrls {
		name, err := ParseGitRemoteUrl(remoteUrl)
		if err != nil {
			t.Error(err)
		}
		if name != debrickedCli {
			t.Error("failed to find correct git repository name:", name)
		}
	}
	badUrl := "ftp://github.com/debricked/cli"
	name, err := ParseGitRemoteUrl(badUrl)
	if err == nil {
		t.Error("failed to assert that bad URL generated error, got:", name)
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
	if !strings.Contains(url, debrickedCli) || !strings.Contains(url, "github.com") {
		t.Error("failed to find correct git remote url:", url)
	}
}

func TestFindRepositoryUrlWithFailure(t *testing.T) {
	repoMock := mockRepository(true, t)
	url, err := FindRepositoryUrl(repoMock)
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
	repoMock := mockRepository(true, t)
	name, err := FindRepositoryName(repoMock, "test/repository")
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if name != "repository" {
		t.Error("failed to find repository name")
	}
}

func TestGetCommit(t *testing.T) {
	repoMock := mockRepository(true, t)
	commit, err := FindCommit(repoMock)
	if err != nil {
		t.Error("failed to assert that an error occurred")
	}
	if commit == nil {
		t.Error("failed to find commit")
	}
}

func TestGetCommitHash(t *testing.T) {
	repoMock := mockRepository(true, t)
	commitHash, err := FindCommitHash(repoMock)
	if err != nil {
		t.Error("failed to assert that an error occurred")
	}
	if len(commitHash) == 0 {
		t.Error("failed assert commit message")
	}
}

func TestGetCommitAuthor(t *testing.T) {
	repoMock := mockRepository(true, t)
	commitHash, err := FindCommitAuthor(repoMock)
	if err != nil {
		t.Error("failed to assert that an error occurred")
	}
	if len(commitHash) == 0 {
		t.Error("failed assert commit message")
	}
}

func TestGetBranch(t *testing.T) {
	repoMock := mockRepository(true, t)
	commitHash, err := FindBranch(repoMock)
	if err != nil {
		t.Error("failed to assert that an error occurred")
	}
	if len(commitHash) == 0 {
		t.Error("failed assert commit message")
	}
}

func mockRepository(withCommit bool, t *testing.T) *git.Repository {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	store := memory.NewStorage()
	r, err := git.Init(store, fs)
	if err != nil {
		t.Fatal("failed to init repository. Error:", err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatal("failed to get worktree. Error:", err)
	}
	if withCommit {
		_, err = w.Commit("Initial commit", &git.CommitOptions{Author: &object.Signature{Name: "author"}})
		if err != nil {
			t.Fatal("failed to create commit. Error:", err)
		}
	}

	return r
}
