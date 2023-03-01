package scan

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/scan"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var repositoryName string
var commitName string
var branchName string
var commitAuthor string
var repositoryUrl string
var integrationName string
var exclusions = file.DefaultExclusions()
var resolve bool

const (
	RepositoryFlag    = "repository"
	CommitFlag        = "commit"
	BranchFlag        = "branch"
	CommitAuthorFlag  = "author"
	RepositoryUrlFlag = "repository-url"
	IntegrationFlag   = "integration"
	ExclusionFlag     = "exclusion"
	ResolveFlag       = "resolve"
)

var scanCmdError error

func NewScanCmd(scanner scan.IScanner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Start a Debricked dependency scan",
		Long: `All supported dependency files will be scanned and analysed.
If the given path contains a git repository all flags but "integration" will be resolved. Otherwise they have to specified.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(&scanner),
	}
	cmd.Flags().StringVarP(&repositoryName, RepositoryFlag, "r", "", "repository name")
	cmd.Flags().StringVarP(&commitName, CommitFlag, "c", "", "commit hash")
	cmd.Flags().StringVarP(&branchName, BranchFlag, "b", "", "branch name")
	cmd.Flags().StringVarP(&commitAuthor, CommitAuthorFlag, "a", "", "commit author")
	cmd.Flags().StringVarP(&repositoryUrl, RepositoryUrlFlag, "u", "", "repository URL")
	cmd.Flags().StringVarP(
		&integrationName,
		IntegrationFlag,
		"i",
		"CLI",
		`name of integration used to trigger scan. For example "GitHub Actions"`,
	)

	fileExclusionExample := filepath.Join("*", "**.lock")
	dirExclusionExample := filepath.Join("**", "node_modules", "**")
	exampleFlags := fmt.Sprintf("-e \"%s\" -e \"%s\"", fileExclusionExample, dirExclusionExample)
	cmd.Flags().StringArrayVarP(
		&exclusions,
		ExclusionFlag,
		"e",
		exclusions,
		`The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Examples: 
$ debricked scan . `+exampleFlags)

	cmd.Flags().BoolVar(&resolve, ResolveFlag, false, `Resolves manifest files that lack lock files. This enables more accurate dependency scanning since the whole dependency tree will be analysed.
For example, if there is a "go.mod" in the target path, its dependencies are going to get resolved onto a lock file, and latter scanned.`)

	viper.MustBindEnv(RepositoryFlag)
	viper.MustBindEnv(CommitFlag)
	viper.MustBindEnv(BranchFlag)
	viper.MustBindEnv(CommitAuthorFlag)
	viper.MustBindEnv(RepositoryUrlFlag)
	viper.MustBindEnv(IntegrationFlag)

	return cmd
}

func RunE(s *scan.IScanner) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		options := scan.DebrickedOptions{
			Path:            path,
			Resolve:         viper.GetBool(ResolveFlag),
			Exclusions:      viper.GetStringSlice(ExclusionFlag),
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

		if scanCmdError == scan.FailPipelineErr {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			return scanCmdError
		} else if scanCmdError != nil {
			return fmt.Errorf("%s %s\n", color.RedString("⨯"), scanCmdError.Error())
		}

		return scanCmdError
	}
}
