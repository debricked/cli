package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"log"
	"path/filepath"
	"regexp"
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
	repository, err := git.PlainOpen(directoryPath)
	if err == nil {
		isSet := func(attribute string) bool { return len(attribute) > 0 }

		if !isSet(repositoryName) {
			repositoryName, err = FindRepositoryName(repository, directoryPath)
			if err != nil {
				log.Println(err.Error())
			}
		}

		head, _ := repository.Head()
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

func checkErrors(obj *MetaObject) error {
	if len(obj.RepositoryName) == 0 {
		return errors.New("failed to find repository name. Please use --repository flag")
	}
	if len(obj.CommitName) == 0 {
		return errors.New("failed to find commit hash. Please use --commit flag")
	}

	return nil
}

func FindRepositoryUrl(repository *git.Repository) (string, error) {
	remoteUrl, err := FindRemoteUrl(repository)
	if err != nil {
		return "", err
	}
	// If remoteUrl starts with "http(s)://" and ends with ".git", use capture group to find repository url.
	var regexes = []string{
		"^(https?:\\/\\/.+)\\.git$",
		"^(https?:\\/\\/.+)$",
	}
	for _, regex := range regexes {
		compiledRegex := regexp.MustCompile(regex)
		matches := compiledRegex.FindStringSubmatch(remoteUrl)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	// If remoteUrl is of the form "git@github.com:organisation/reponame.git",
	// use capture groups to construct repository url
	const gitUrlRegex = "git@(.+):[0-9]*\\/?(.+)\\.git$"
	compiledRegex := regexp.MustCompile(gitUrlRegex)
	matches := compiledRegex.FindStringSubmatch(remoteUrl)
	if len(matches) > 2 {
		domain := matches[1]
		uri := matches[2]
		url := fmt.Sprintf("https://%s/%s", domain, uri)

		return url, nil
	}

	return "", errors.New("failed to find repository URL")
}

// FindRemoteUrl returns first remote URL found in the repository
func FindRemoteUrl(repository *git.Repository) (string, error) {
	var err error = nil
	remoteURL := ""
	remotes, _ := repository.Remotes()
	for _, remote := range remotes {
		remoteLinks := remote.Config().URLs
		for _, remoteLink := range remoteLinks {
			remoteURL = remoteLink
			break
		}
		if remoteURL != "" {
			break
		}
	}

	if remoteURL == "" {
		err = errors.New("failed to find repository remote URL")
	}
	return remoteURL, err
}

func FindRepositoryName(repository *git.Repository, directoryPath string) (string, error) {
	absolutePath, _ := filepath.Abs(directoryPath)
	repositoryName := filepath.Base(absolutePath)
	gitRemoteUrl, err := FindRemoteUrl(repository)
	if err != nil {
		return repositoryName, err
	}

	return ParseGitRemoteUrl(gitRemoteUrl)
}

func ParseGitRemoteUrl(gitRemoteUrl string) (string, error) {
	const httpsUrlRegex = "https?:\\/\\/.+\\.[a-z0-9]+\\/(.+)\\.git$"
	const urlRegex = "https?:\\/\\/.+\\.[a-z0-9]+\\/(.+)$"
	const gitUrlRegex = "^.*:[0-9]*\\/*(.+)\\.git$"
	// 1. Try to match against https git URL
	// 2. Try to match against ssh git URL
	regexes := []string{httpsUrlRegex, urlRegex, gitUrlRegex}
	for _, regex := range regexes {
		compiledRegex := regexp.MustCompile(regex)
		matches := compiledRegex.FindStringSubmatch(gitRemoteUrl)
		if matches != nil && len(matches) > 1 {
			return matches[1], nil
		}
	}

	return gitRemoteUrl, errors.New("failed to parse git remote URL. git/https regular expressions had no matches")
}
