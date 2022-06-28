package scan

import (
	"debricked/pkg/client"
	"debricked/pkg/file"
	"debricked/pkg/git"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var debClient *client.DebClient
var finder *file.Finder

var repositoryName string
var commitName string
var branchName string
var commitAuthor string
var repositoryUrl string
var integrationName string
var exclusions []string

func NewScanCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	finder, _ = file.NewFinder(debClient)
	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Start a Debricked dependency scan",
		Long: `All supported dependency files will be scanned and analysed.
If the given path contains a git repository all flags but "integration" will be resolved. Otherwise they have to specified.`,
		Args: validateArgs,
		RunE: run,
	}

	cmd.Flags().StringVarP(&repositoryName, "repository", "r", "", "repository name")
	cmd.Flags().StringVarP(&commitName, "commit", "c", "", "commit hash")
	cmd.Flags().StringVarP(&branchName, "branch", "b", "", "branch name")
	cmd.Flags().StringVarP(&commitAuthor, "author", "a", "", "commit author")
	cmd.Flags().StringVarP(&repositoryUrl, "repository-url", "u", "", "repository URL")
	cmd.Flags().StringVarP(&integrationName, "integration", "i", "CLI", `name of integration used to trigger scan. For example "GitHub Actions"`)
	cmd.Flags().StringArrayVarP(&exclusions, "exclude", "e", exclusions, `The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Examples: 
$ debricked scan . -e "*/**.lock" -e "**/node_modules/**" 
$ debricked scan . -e "*\**.exe" -e "**\node_modules\**" 
`)

	return cmd
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires directory path")
	}
	if isValidFilepath(args[0]) {
		return nil
	}
	return fmt.Errorf("invalid directory path specified: %s", args[0])
}

func run(_ *cobra.Command, args []string) error {
	directoryPath := args[0]
	gitMetaObject, err := git.NewMetaObject(directoryPath, repositoryName, commitName, branchName, commitAuthor, repositoryUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s\n", color.RedString("⨯"), err.Error()))
	}
	err = scan(directoryPath, gitMetaObject, exclusions)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s\n", color.RedString("⨯"), err.Error()))
	}

	return nil
}

func isValidFilepath(path string) bool {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}

	return true
}

func scan(directoryPath string, gitMetaObject *git.MetaObject, exclusions []string) error {
	fileGroups, err := finder.GetGroups(directoryPath, exclusions)
	if err != nil {
		return err
	}

	batch := newUploadBatch(fileGroups, gitMetaObject)
	batch.upload()
	err = batch.conclude()
	if err != nil {
		return err
	}

	result, err := batch.wait()
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
