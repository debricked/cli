package testdata

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
)

func SetUpCiEnv(env map[string]string) error {
	for variable, value := range env {
		err := os.Setenv(variable, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetUpGitRepository(includeCommit bool) (string, error) {
	cwd, _ := os.Getwd()
	repoDir := cwd + "/testdata/"
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		return cwd, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return cwd, err
	}
	if includeCommit {
		_, err = w.Commit("Initial commit", &git.CommitOptions{Author: &object.Signature{Name: "author"}})
		if err != nil {
			return cwd, err
		}
	}

	err = os.Chdir(repoDir)
	return cwd, err
}

func TearDownGitRepository(dir string) error {
	cwd, _ := os.Getwd()
	repoDir := cwd + "/.git/"
	_, err := git.PlainOpen(repoDir)
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	return os.RemoveAll(repoDir)
}

func ResetEnv(ciEnv map[string]string) error {
	for _, variable := range ciEnv {
		err := os.Unsetenv(variable)
		if err != nil {
			return err
		}
	}
	return nil
}
