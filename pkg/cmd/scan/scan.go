package scan

import (
	"debricked/pkg/ci"
	"debricked/pkg/client"
	"debricked/pkg/scan"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var repositoryName string
var commitName string
var branchName string
var commitAuthor string
var repositoryUrl string
var integrationName string
var exclusions []string

const (
	RepositoryFlag    = "repository"
	CommitFlag        = "commit"
	BranchFlag        = "branch"
	CommitAuthorFlag  = "author"
	RepositoryUrlFlag = "repository-url"
	IntegrationFlag   = "integration"
	ExclusionsFlag    = "exclusions"
)

var scanCmdError error

func NewScanCmd(c *client.IDebClient) *cobra.Command {
	var ciService ci.IService
	ciService = ci.NewService(nil)

	var s scan.IScanner
	s, scanCmdError = scan.NewDebrickedScanner(c, ciService)

	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Start a Debricked dependency scan",
		Long: `All supported dependency files will be scanned and analysed.
If the given path contains a git repository all flags but "integration" will be resolved. Otherwise they have to specified.`,
		Args: ValidateArgs,
		RunE: RunE(&s),
	}
	cmd.Flags().StringVarP(&repositoryName, RepositoryFlag, "r", "", "repository name")
	cmd.Flags().StringVarP(&commitName, CommitFlag, "c", "", "commit hash")
	cmd.Flags().StringVarP(&branchName, BranchFlag, "b", "", "branch name")
	cmd.Flags().StringVarP(&commitAuthor, CommitAuthorFlag, "a", "", "commit author")
	cmd.Flags().StringVarP(&repositoryUrl, RepositoryUrlFlag, "u", "", "repository URL")
	cmd.Flags().StringVarP(&integrationName, IntegrationFlag, "i", "CLI", `name of integration used to trigger scan. For example "GitHub Actions"`)
	cmd.Flags().StringArrayVarP(&exclusions, ExclusionsFlag, "e", exclusions, `The following terms are supported to exclude paths:
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
	viper.MustBindEnv(RepositoryFlag)
	viper.MustBindEnv(CommitFlag)
	viper.MustBindEnv(BranchFlag)
	viper.MustBindEnv(CommitAuthorFlag)
	viper.MustBindEnv(RepositoryUrlFlag)
	viper.MustBindEnv(IntegrationFlag)
	viper.MustBindEnv(ExclusionsFlag)

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

func RunE(s *scan.IScanner) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		directoryPath := args[0]
		options := scan.DebrickedOptions{
			DirectoryPath:   directoryPath,
			Exclusions:      viper.GetStringSlice(ExclusionsFlag),
			RepositoryName:  viper.GetString(RepositoryFlag),
			CommitName:      viper.GetString(CommitFlag),
			BranchName:      viper.GetString(BranchFlag),
			CommitAuthor:    viper.GetString(CommitAuthorFlag),
			RepositoryUrl:   viper.GetString(RepositoryUrlFlag),
			IntegrationName: viper.GetString(IntegrationFlag),
		}
		if s != nil {
			scanCmdError = (*s).Scan(options)
		} else {
			scanCmdError = errors.New("scanner was nil")
		}

		if scanCmdError != nil {
			return errors.New(fmt.Sprintf("%s %s\n", color.RedString("⨯"), scanCmdError.Error()))
		}

		return scanCmdError
	}
}

func ValidateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires directory path")
	}
	if isValidFilepath(args[0]) {
		return nil
	}
	return fmt.Errorf("invalid directory path specified: %s", args[0])
}

func isValidFilepath(path string) bool {
	_, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	return true
}
