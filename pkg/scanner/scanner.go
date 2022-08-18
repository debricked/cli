package scanner

import (
	"debricked/pkg/client"
	"debricked/pkg/file"
	"debricked/pkg/git"
	"debricked/pkg/uploader"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
)

type Scanner interface {
	Scan(o Options) error
}

type Options interface{}

type debrickedScanner struct {
	client   *client.Client
	finder   *file.Finder
	uploader *uploader.Uploader
}

type DebrickedOptions struct {
	DirectoryPath   string
	Exclusions      []string
	RepositoryName  string
	CommitName      string
	BranchName      string
	CommitAuthor    string
	RepositoryUrl   string
	IntegrationName string
}

func NewDebrickedScanner(c *client.Client) (*debrickedScanner, error) {
	finder, err := file.NewFinder(*c)
	if err != nil {
		return nil, newInitError(err)
	}
	var u uploader.Uploader
	u, err = uploader.NewDebrickedUploader(c)

	if err != nil {
		return nil, newInitError(err)
	}

	return &debrickedScanner{
		c,
		finder,
		&u,
	}, nil
}

func (dScanner *debrickedScanner) Scan(o Options) error {
	dOptions := o.(DebrickedOptions)
	gitMetaObject, err := git.NewMetaObject(
		dOptions.DirectoryPath,
		dOptions.RepositoryName,
		dOptions.CommitName,
		dOptions.BranchName,
		dOptions.CommitAuthor,
		dOptions.RepositoryUrl,
	)
	if err != nil {
		return err
	}

	fileGroups, err := dScanner.finder.GetGroups(dOptions.DirectoryPath, dOptions.Exclusions)
	if err != nil {
		return err
	}

	uploaderOptions := uploader.DebrickedOptions{FileGroups: fileGroups, GitMetaObject: *gitMetaObject, IntegrationsName: dOptions.IntegrationName}
	result, err := (*dScanner.uploader).Upload(uploaderOptions)
	if err != nil {
		return err
	}

	fmt.Printf("\n%d vulnerabilities found\n", result.VulnerabilitiesFound)
	fmt.Println("")
	failPipeline := false
	for _, rule := range result.AutomationRules {
		rule.Print(os.Stdout)
		failPipeline = failPipeline || rule.FailPipeline()
	}
	fmt.Printf("For full details, visit: %s\n\n", color.BlueString(result.DetailsUrl))
	if failPipeline {
		return errors.New("")
	}

	return nil
}

func newInitError(err error) error {
	return errors.New("failed to initialize the uploader due to: " + err.Error())
}
