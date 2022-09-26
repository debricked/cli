package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"path/filepath"
	"regexp"
)

var RepositoryNameError = errors.New("failed to find repository name")
var CommitNameError = errors.New("failed to find commit hash")

func checkErrors(obj *MetaObject) error {
	if len(obj.RepositoryName) == 0 {
		return RepositoryNameError
	}
	if len(obj.CommitName) == 0 {
		return CommitNameError
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

func FindRepository(directoryPath string) (*git.Repository, error) {
	return git.PlainOpen(directoryPath)
}

func FindBranch(repository *git.Repository) (string, error) {
	head, err := repository.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}

func FindCommit(repository *git.Repository) (*object.Commit, error) {
	head, err := repository.Head()
	if err != nil {
		return nil, err
	}
	commitObject, err := repository.CommitObject(head.Hash())

	return commitObject, err
}

func FindCommitAuthor(repository *git.Repository) (string, error) {
	c, err := FindCommit(repository)
	if err != nil {
		return "", err
	}
	return c.Author.String(), nil
}

func FindCommitHash(repository *git.Repository) (string, error) {
	c, err := FindCommit(repository)
	if err != nil {
		return "", err
	}
	return c.Hash.String(), nil
}
