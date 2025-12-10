package scan

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
var generateCommitName bool
var debug bool
var exclusions = file.Exclusions()
var inclusions []string
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
var sbom string
var sbomOutput string
var tagCommitAsRelease bool
var experimental bool

const (
	BranchFlag                      = "branch"
	CallGraphFlag                   = "callgraph"
	CallGraphGenerateTimeoutFlag    = "callgraph-generate-timeout"
	CallGraphUploadTimeoutFlag      = "callgraph-upload-timeout"
	CommitFlag                      = "commit"
	CommitAuthorFlag                = "author"
	DebugFlag                       = "debug"
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
	SBOMFlag                        = "sbom"
	SBOMOutputFlag                  = "sbom-output"
	TagCommitAsReleaseFlag          = "tag-commit-as-release"
	TagCommitAsReleaseEnv           = "TAG_COMMIT_AS_RELEASE"
	ExperimentalFlag                = "experimental"
	GenerateCommitNameFlag          = "generate-commit-name"
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
			return RunE(&scanner)(cmd, args)
		},
	}

	cmd.Flags().StringVarP(&repositoryName, RepositoryFlag, "r", "", "repository name")
	cmd.Flags().StringVarP(&commitName, CommitFlag, "c", "", "commit hash")
	cmd.Flags().BoolVar(&generateCommitName, GenerateCommitNameFlag, false, "auto-generate a commit name if no commit hash is found (in -c or env)")
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
$ debricked scan . --inclusion '**/node_modules/**'`)
	regenerateDoc := strings.Join(
		[]string{
			"Toggle regeneration of already existing lock files between 3 modes:\n",
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
			"Toggle version hinting, i.e using manifest versions to help manifestless resolution.\n",
			"\nExample:\n$ debricked scan . --version-hint=false",
		}, "\n")
	cmd.Flags().BoolVar(&versionHint, VersionHintFlag, true, versionHintDoc)
	experimentalFlagDoc := strings.Join(
		[]string{
			"This flag allows inclusion of repository matches",
			"\nExample:\n$ debricked scan . --experimental=false",
		}, "\n")
	cmd.Flags().BoolVar(&experimental, ExperimentalFlag, false, experimentalFlagDoc)
	verboseDoc := strings.Join(
		[]string{
			"This flag allows you to reduce error output for resolution.",
			"\nExample:\n$ debricked resolve --verbose=false",
		}, "\n")
	cmd.Flags().BoolVar(&verbose, VerboseFlag, true, verboseDoc)
	cmd.Flags().BoolVar(&debug, DebugFlag, false, "write all debug output to stderr")
	cmd.Flags().BoolVarP(&passOnDowntime, PassOnTimeOut, "p", false, "pass scan if there is a service access timeout")
	cmd.Flags().BoolVar(&noResolve, NoResolveFlag, false, `disables resolution of manifest files that lack lock files. Resolving manifest files enables more accurate dependency scanning since the whole dependency tree will be analysed.
For example, if there is a "go.mod" in the target path, its dependencies are going to get resolved onto a lock file, and latter scanned.`)
	cmd.Flags().BoolVar(&noFingerprint, NoFingerprintFlag, false, "Toggle fingerprinting for undeclared component identification. Can be run as a standalone command [fingerprint] with more granular options.")
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
	cmd.Flags().StringVar(&sbom, SBOMFlag, "", `Toggle generating and downloading SBOM report after scan completion of specified format.
Supported formats are: 'CycloneDX', 'SPDX'
Leaving the field empty results in no SBOM generation.`,
	)
	cmd.Flags().StringVar(&sbomOutput, SBOMOutputFlag, "", `Set output path of downloaded SBOM report (if sbom is toggled)`)
	cmd.Flags().BoolVar(
		&tagCommitAsRelease,
		TagCommitAsReleaseFlag,
		false,
		"Set to true to tag commit as a release. This will store the scan data indefinitely. Enterprise is required for this flag. Please visit https://debricked.com/pricing/ for more info. Can be overridden by "+TagCommitAsReleaseEnv+" environment variable.",
	)

	viper.MustBindEnv(RepositoryFlag)
	viper.MustBindEnv(CommitFlag)
	viper.MustBindEnv(BranchFlag)
	viper.MustBindEnv(CommitAuthorFlag)
	viper.MustBindEnv(RepositoryUrlFlag)
	viper.MustBindEnv(IntegrationFlag)
	viper.MustBindEnv(PassOnTimeOut)
	viper.MustBindEnv(NpmPreferredFlag)
	viper.MustBindEnv(SBOMFlag)
	viper.MustBindEnv(SBOMOutputFlag)
	viper.MustBindEnv(TagCommitAsReleaseFlag)

	// Hide experimental flag
	err := cmd.Flags().MarkHidden(ExperimentalFlag)
	if err != nil { // This should not be reachable
		fmt.Println("Trying to hide non-existing flag")
	}

	return cmd
}

func RunE(s *scan.IScanner) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Scanner started...")

		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		tagCommitAsRelease := false
		tagCommitAsReleaseEnv := os.Getenv(TagCommitAsReleaseEnv)
		if tagCommitAsReleaseEnv != "" {
			var err error
			tagCommitAsRelease, err = strconv.ParseBool(tagCommitAsReleaseEnv)

			if err != nil {
				return errors.Join(errors.New("failed to convert "+TagCommitAsReleaseEnv+" to boolean"), err)
			}
		} else {
			tagCommitAsRelease = viper.GetBool(TagCommitAsReleaseFlag)
		}

		options := scan.DebrickedOptions{
			Path:                        path,
			Resolve:                     !viper.GetBool(NoResolveFlag),
			Fingerprint:                 !viper.GetBool(NoFingerprintFlag),
			SBOM:                        viper.GetString(SBOMFlag),
			SBOMOutput:                  viper.GetString(SBOMOutputFlag),
			Exclusions:                  viper.GetStringSlice(ExclusionFlag),
			Inclusions:                  viper.GetStringSlice(InclusionFlag),
			Verbose:                     viper.GetBool(VerboseFlag),
			Debug:                       viper.GetBool(DebugFlag),
			Regenerate:                  viper.GetInt(RegenerateFlag),
			VersionHint:                 viper.GetBool(VersionHintFlag),
			RepositoryName:              viper.GetString(RepositoryFlag),
			CommitName:                  viper.GetString(CommitFlag),
			GenerateCommitName:          viper.GetBool(GenerateCommitNameFlag),
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
			TagCommitAsRelease:          tagCommitAsRelease,
			Experimental:                viper.GetBool(ExperimentalFlag),
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
			return fmt.Errorf("%s\n", color.RedString("⨯")+" "+scanCmdError.Error())
		}

		return scanCmdError
	}
}
