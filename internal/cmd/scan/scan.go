package scan

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/scan"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var branchName string
var callgraph bool
var callgraphGenerateTimeout int
var callgraphUploadTimeout int
var commitAuthor string
var commitName string
var exclusions = file.Exclusions()
var inclusions = file.Exclusions()
var integrationName string
var jsonFilePath string
var minFingerprintContentLength int
var noFingerprint bool
var noResolve bool
var npmPreferred bool
var passOnDowntime bool
var regenerate int
var repositoryName string
var repositoryUrl string
var verbose bool
var versionHint bool

const (
	BranchFlag                      = "branch"
	CallGraphFlag                   = "callgraph"
	CallGraphGenerateTimeoutFlag    = "callgraph-generate-timeout"
	CallGraphUploadTimeoutFlag      = "callgraph-upload-timeout"
	CommitFlag                      = "commit"
	CommitAuthorFlag                = "author"
	ExclusionFlag                   = "exclusion"
	IntegrationFlag                 = "integration"
	InclusionFlag                   = "inclusion"
	JsonFilePathFlag                = "json-path"
	MinFingerprintContentLengthFlag = "min-fingerprint-content-length"
	NoResolveFlag                   = "no-resolve"
	NoFingerprintFlag               = "no-fingerprint"
	NpmPreferredFlag                = "prefer-npm"
	PassOnTimeOut                   = "pass-on-timeout"
	RegenerateFlag                  = "regenerate"
	RepositoryFlag                  = "repository"
	RepositoryUrlFlag               = "repository-url"
	VerboseFlag                     = "verbose"
	VersionHintFlag                 = "version-hint"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(viper.GetString(RepositoryFlag)) > 0 {
				if strings.ToLower(viper.GetString(RepositoryFlag))[0] < 'm' && !cmd.Flags().Changed(NoFingerprintFlag) {
					viper.Set(NoFingerprintFlag, false)
				} // Temporary addition for rolling release of fingerprinting enabled by default
			}

			return RunE(&scanner)(cmd, args)
		},
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
	cmd.Flags().StringVarP(&jsonFilePath, JsonFilePathFlag, "j", "", "write upload result as json to provided path")
	fileExclusionExample := filepath.Join("'*", "**.lock'")
	dirExclusionExample := filepath.Join("'**", "node_modules", "**'")
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

Exclude flags could alternatively be set using DEBRICKED_EXCLUSIONS="path1,path2,path3".

Examples: 
$ debricked scan . `+exampleFlags)
	cmd.Flags().StringArrayVar(
		&inclusions,
		InclusionFlag,
		inclusions,
		`Forces inclusion of specified terms, see exclusion flag for more information on supported terms.
Examples: 
$ debricked scan . --include '**/node_modules/**'`)
	regenerateDoc := strings.Join(
		[]string{
			"Toggles regeneration of already existing lock files between 3 modes:\n",
			"Force Regeneration Level | Meaning",
			"------------------------ | -------",
			"0 (default)              | No regeneration",
			"1                        | Regenerates existing non package manager native Debricked lock files",
			"2                        | Regenerates all existing lock files",
			"\nExample:\n$ debricked resolve . --regenerate=1",
		}, "\n")
	cmd.Flags().IntVar(&regenerate, RegenerateFlag, 0, regenerateDoc)
	versionHintDoc := strings.Join(
		[]string{
			"Toggles version hinting, i.e using manifest versions to help manifestless resolution.\n",
			"\nExample:\n$ debricked scan . --version-hint=false",
		}, "\n")
	cmd.Flags().BoolVar(&versionHint, VersionHintFlag, true, versionHintDoc)
	verboseDoc := strings.Join(
		[]string{
			"This flag allows you to reduce error output for resolution.",
			"\nExample:\n$ debricked resolve --verbose=false",
		}, "\n")
	cmd.Flags().BoolVar(&verbose, VerboseFlag, true, verboseDoc)
	cmd.Flags().BoolVarP(&passOnDowntime, PassOnTimeOut, "p", false, "pass scan if there is a service access timeout")
	cmd.Flags().BoolVar(&noResolve, NoResolveFlag, false, `disables resolution of manifest files that lack lock files. Resolving manifest files enables more accurate dependency scanning since the whole dependency tree will be analysed.
For example, if there is a "go.mod" in the target path, its dependencies are going to get resolved onto a lock file, and latter scanned.`)
	cmd.Flags().BoolVar(&noFingerprint, NoFingerprintFlag, true, "toggles fingerprinting for undeclared component identification. Can be run as a standalone command [fingerprint] with more granular options.")
	cmd.Flags().BoolVar(&callgraph, CallGraphFlag, false, `Enables call graph generation during scan.`)
	cmd.Flags().IntVar(&callgraphUploadTimeout, CallGraphUploadTimeoutFlag, 10*60, "Set a timeout (in seconds) on call graph upload.")
	cmd.Flags().IntVar(&callgraphGenerateTimeout, CallGraphGenerateTimeoutFlag, 60*60, "Set a timeout (in seconds) on call graph generation.")
	cmd.Flags().IntVar(&minFingerprintContentLength, MinFingerprintContentLengthFlag, 0, "Set minimum content length (in bytes) for files to fingerprint.")
	npmPreferredDoc := strings.Join(
		[]string{
			"This flag allows you to select which package manager will be used as a resolver: Yarn (default) or NPM.",
			"Example: debricked resolve --prefer-npm",
		}, "\n")
	cmd.Flags().BoolP(NpmPreferredFlag, "", npmPreferred, npmPreferredDoc)

	viper.MustBindEnv(RepositoryFlag)
	viper.MustBindEnv(CommitFlag)
	viper.MustBindEnv(BranchFlag)
	viper.MustBindEnv(CommitAuthorFlag)
	viper.MustBindEnv(RepositoryUrlFlag)
	viper.MustBindEnv(IntegrationFlag)
	viper.MustBindEnv(PassOnTimeOut)
	viper.MustBindEnv(NpmPreferredFlag)

	return cmd
}

func RunE(s *scan.IScanner) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		options := scan.DebrickedOptions{
			Path:                        path,
			Resolve:                     !viper.GetBool(NoResolveFlag),
			Fingerprint:                 !viper.GetBool(NoFingerprintFlag),
			Exclusions:                  viper.GetStringSlice(ExclusionFlag),
			Verbose:                     viper.GetBool(VerboseFlag),
			Regenerate:                  viper.GetInt(RegenerateFlag),
			VersionHint:                 viper.GetBool(VersionHintFlag),
			RepositoryName:              viper.GetString(RepositoryFlag),
			CommitName:                  viper.GetString(CommitFlag),
			BranchName:                  viper.GetString(BranchFlag),
			CommitAuthor:                viper.GetString(CommitAuthorFlag),
			RepositoryUrl:               viper.GetString(RepositoryUrlFlag),
			IntegrationName:             viper.GetString(IntegrationFlag),
			JsonFilePath:                viper.GetString(JsonFilePathFlag),
			NpmPreferred:                viper.GetBool(NpmPreferredFlag),
			PassOnTimeOut:               viper.GetBool(PassOnTimeOut),
			CallGraph:                   viper.GetBool(CallGraphFlag),
			CallGraphUploadTimeout:      viper.GetInt(CallGraphUploadTimeoutFlag),
			CallGraphGenerateTimeout:    viper.GetInt(CallGraphGenerateTimeoutFlag),
			MinFingerprintContentLength: viper.GetInt(MinFingerprintContentLengthFlag),
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
