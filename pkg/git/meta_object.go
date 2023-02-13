package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"log"
)

type MetaObject struct {
	RepositoryName    string
	CommitName        string
	RepositoryUrl     string
	BranchName        string
	DefaultBranchName string
	Author            string
}

// NewMetaObject returns MetaObject based on git repository existing on path. Otherwise, inputted arguments are used
func NewMetaObject(path string, repositoryName string, commit string, branchName string, commitAuthor string, url string) (*MetaObject, error) {
	repository, repoErr := FindRepository(path)
	if repoErr == nil {
		setRepositoryName(&repositoryName, repository, path)
		head, headErr := repository.Head()
		if headErr == nil {
			setCommit(&commit, repository, *head)
			setCommitAuthor(&commitAuthor, repository, *head)
			setBranch(&branchName, *head)
			setRepositoryUrl(&url, repository)
		} else {
			log.Println(headErr.Error())
		}
	}

	obj := &MetaObject{
		RepositoryName:    repositoryName,
		CommitName:        commit,
		RepositoryUrl:     url,
		BranchName:        branchName,
		DefaultBranchName: branchName,
		Author:            commitAuthor,
	}

	return obj, checkErrors(obj)
}

func isSet(attribute string) bool {
	return len(attribute) > 0
}

func setRepositoryName(name *string, repository *git.Repository, path string) {
	var err error
	if !isSet(*name) {
		*name, err = FindRepositoryName(repository, path)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func setCommit(commit *string, repository *git.Repository, head plumbing.Reference) {
	commitObject, _ := repository.CommitObject(head.Hash())
	if !isSet(*commit) {
		*commit = commitObject.Hash.String()
	}
}

func setCommitAuthor(commitAuthor *string, repository *git.Repository, head plumbing.Reference) {
	commitObject, _ := repository.CommitObject(head.Hash())
	if !isSet(*commitAuthor) {
		*commitAuthor = commitObject.Author.String()
	}
}

func setBranch(branch *string, head plumbing.Reference) {
	if !isSet(*branch) {
		*branch = head.Name().Short()
	}
}

func setRepositoryUrl(url *string, repository *git.Repository) {
	var err error
	if !isSet(*url) {
		*url, err = FindRepositoryUrl(repository)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
