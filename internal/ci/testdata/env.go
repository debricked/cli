package testdata

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func SetUpCiEnv(t *testing.T, env map[string]string) {
	for variable, value := range env {
		err := os.Setenv(variable, value)
		if err != nil {
			t.Fatal("failed to set up Ci env. Err: ", err)
		}
	}
}

func SetUpGitRepository(t *testing.T, includeCommit bool) string {
	cwd, _ := os.Getwd()
	repoDir := cwd + "/testdata/"
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if includeCommit {
		tempFilepath := repoDir + "test.txt"
		f, err := os.Create(tempFilepath)
		if err != nil {
			t.Fatal("failed to create file. Error:", err)
		}
		err = f.Close()
		if err != nil {
			t.Fatal("failed to close created file. Error:", err)
		}
		_, err = w.Add("test.txt")
		if err != nil {
			t.Fatal("failed to add file to worktree. Error:", err)
		}
		_, err = w.Commit("Initial commit", &git.CommitOptions{Author: &object.Signature{Name: "author"}})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = os.Chdir(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	return cwd
}

func TearDownGitRepository(dir string, t *testing.T) {
	cwd, _ := os.Getwd()
	repoDir := cwd
	_, err := git.PlainOpen(repoDir)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll(repoDir)
	if err != nil {
		t.Fatal(err)
	}
}

func ResetEnv(t *testing.T, ciEnv map[string]string) {
	for _, variable := range ciEnv {
		UnsetEnvVar(t, variable)
	}
}

func UnsetEnvVar(t *testing.T, envVar string) {
	err := os.Unsetenv(envVar)
	if err != nil {
		t.Fatal(err)
	}
}
