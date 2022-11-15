package git

import (
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

// NewMetaObject returns MetaObject based on git repository existing on directoryPath. Otherwise, inputted arguments are used
func NewMetaObject(directoryPath string, repositoryName string, commit string, branchName string, commitAuthor string, url string) (*MetaObject, error) {
	repository, err := FindRepository(directoryPath)
	if err == nil {
		isSet := func(attribute string) bool { return len(attribute) > 0 }

		if !isSet(repositoryName) {
			repositoryName, err = FindRepositoryName(repository, directoryPath)
			if err != nil {
				log.Println(err.Error())
			}
		}

		head, err := repository.Head()
		if err == nil {
			commitObject, _ := repository.CommitObject(head.Hash())
			if !isSet(commit) {
				commit = commitObject.Hash.String()
			}
			if !isSet(commitAuthor) {
				commitAuthor = commitObject.Author.String()
			}

			if !isSet(branchName) {
				branchName = head.Name().Short()
			}

			if !isSet(url) {
				url, err = FindRepositoryUrl(repository)
				if err != nil {
					log.Println(err.Error())
				}
			}
		} else {
			log.Println(err.Error())
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
