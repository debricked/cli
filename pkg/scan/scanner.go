package scan

import (
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/ci"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"github.com/debricked/cli/pkg/tui"
	"github.com/debricked/cli/pkg/upload"
	"github.com/fatih/color"
	"os"
	"path/filepath"
)

var (
	BadOptsErr      = errors.New("failed to type case IOptions")
	FailPipelineErr = errors.New("")
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
	Path            string
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
	u, err = upload.NewUploader(*c)

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

	MapEnvToOptions(&dOptions, e)

	if err := SetWorkingDirectory(&dOptions); err != nil {
		return err
	}

	gitMetaObject, err := git.NewMetaObject(
		dOptions.Path,
		dOptions.RepositoryName,
		dOptions.CommitName,
		dOptions.BranchName,
		dOptions.CommitAuthor,
		dOptions.RepositoryUrl,
	)
	if err != nil {
		return err
	}

	fileGroups, err := dScanner.finder.GetGroups(dOptions.Path, dOptions.Exclusions, false, file.StrictAll)
	if err != nil {
		return err
	}

	uploaderOptions := upload.DebrickedOptions{FileGroups: fileGroups, GitMetaObject: *gitMetaObject, IntegrationsName: dOptions.IntegrationName}
	result, err := (*dScanner.uploader).Upload(uploaderOptions)
	if err != nil {
		return err
	}

	if result == nil {
		fmt.Println("Progress polling terminated due to long scan times. Please try again later")

		return nil
	}

	fmt.Printf("\n%d vulnerabilities found\n", result.VulnerabilitiesFound)
	fmt.Println("")
	failPipeline := false
	for _, rule := range result.AutomationRules {
		tui.NewRuleCard(os.Stdout, rule).Render()
		failPipeline = failPipeline || (rule.Triggered && rule.FailPipeline())
	}
	fmt.Printf("For full details, visit: %s\n\n", color.BlueString(result.DetailsUrl))
	if failPipeline {
		return FailPipelineErr
	}

	return nil
}

// SetWorkingDirectory sets working directory in accordance with the path option
func SetWorkingDirectory(d *DebrickedOptions) error {
	absPath, _ := filepath.Abs(d.Path)
	err := os.Chdir(absPath)
	if err != nil {
		return err
	}
	d.Path = ""
	fmt.Printf("Working directory: %s\n", absPath)

	return nil
}

func MapEnvToOptions(o *DebrickedOptions, env env.Env) {
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
	if o.IntegrationName == "CLI" {
		if len(env.Integration) != 0 {
			o.IntegrationName = env.Integration
		}
	}
	if len(o.Path) == 0 && len(env.Filepath) > 0 {
		o.Path = env.Filepath
	}
}

func newInitError(err error) error {
	return errors.New("failed to initialize the uploader due to: " + err.Error())
}
