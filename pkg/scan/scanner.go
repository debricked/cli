package scan

import (
	"debricked/pkg/ci"
	"debricked/pkg/ci/env"
	"debricked/pkg/client"
	"debricked/pkg/file"
	"debricked/pkg/git"
	"debricked/pkg/upload"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
)

var (
	BadOptsErr = errors.New("failed to type case IOptions")
)

type IScanner interface {
	Scan(o IOptions) error
}

type IOptions interface{}

type DebrickedScanner struct {
	client    *client.IDebClient
	finder    *file.Finder
	uploader  *upload.IUploader
	ciService ci.IService
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

func NewDebrickedScanner(c *client.IDebClient, ciService ci.IService) (*DebrickedScanner, error) {
	finder, err := file.NewFinder(*c)
	if err != nil {
		return nil, newInitError(err)
	}
	var u upload.IUploader
	u, err = upload.NewUploader(c)

	if err != nil {
		return nil, newInitError(err)
	}

	return &DebrickedScanner{
		c,
		finder,
		&u,
		ciService,
	}, nil
}

func (dScanner *DebrickedScanner) Scan(o IOptions) error {
	dOptions, ok := o.(DebrickedOptions)
	if !ok {
		return BadOptsErr
	}
	e, _ := dScanner.ciService.Find()
	dScanner.mapEnvToOptions(&dOptions, e)

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

	uploaderOptions := upload.DebrickedOptions{FileGroups: fileGroups, GitMetaObject: *gitMetaObject, IntegrationsName: dOptions.IntegrationName}
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

func (dScanner *DebrickedScanner) mapEnvToOptions(o *DebrickedOptions, env env.Env) {
	if len(o.RepositoryName) == 0 {
		o.RepositoryName = env.Repository
	}
	if len(o.CommitName) == 0 {
		o.CommitName = env.Commit
	}
	if len(o.BranchName) == 0 {
		o.BranchName = env.Branch
	}
	if len(o.CommitAuthor) == 0 {
		o.CommitAuthor = env.Author
	}
	if len(o.RepositoryUrl) == 0 {
		o.RepositoryUrl = env.RepositoryUrl
	}
	if len(o.IntegrationName) == 0 {
		o.IntegrationName = env.Integration
	}
	if len(o.DirectoryPath) == 0 {
		o.DirectoryPath = env.Filepath
	}
}

func newInitError(err error) error {
	return errors.New("failed to initialize the uploader due to: " + err.Error())
}
