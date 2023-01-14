package scan

import (
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/ci"
	"github.com/debricked/cli/pkg/docker"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"github.com/debricked/cli/pkg/upload"
	"github.com/fatih/color"
	"os"
)

type DockerScanner struct {
	Uploader    upload.IUploader
	SbomFactory docker.ISbomFactory
	CiService   ci.IService
}

func (dScanner *DockerScanner) Scan(o IOptions) error {
	dOptions, ok := o.(DockerOptions)
	if !ok {
		return BadOptsErr
	}

	img := docker.MakeImage(dOptions.Image)
	defer func() {
		_ = os.Remove(docker.SbomFileName)
	}()
	sbom, err := dScanner.SbomFactory.Make(&img)
	if err != nil {
		return err
	}

	e, _ := dScanner.CiService.Find()
	if len(e.Integration) == 0 {
		e.Integration = "cli"
	}

	fileGroups := file.Groups{}
	fileGroups.Add(file.Group{FilePath: sbom.Name()})
	result, err := dScanner.Uploader.Upload(upload.DebrickedOptions{
		FileGroups: fileGroups,
		GitMetaObject: git.MetaObject{
			RepositoryName: fmt.Sprintf("🐳 %s", img.Repository),
			CommitName:     img.Digest,
			BranchName:     img.Tag,
		},
		IntegrationsName: e.Integration,
	})
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
