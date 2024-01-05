package resolve

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/resolution"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	exclusions   = file.Exclusions()
	verbose      bool
	npmPreferred bool
	regenerate   int
	strictness   int
)

const (
	ExclusionFlag        = "exclusion"
	VerboseFlag          = "verbose"
	NpmPreferredFlag     = "prefer-npm"
	RegenerateFlag       = "regenerate"
	ResolutionStrictFlag = "resolution-strictness"
)

func NewResolveCmd(resolver resolution.IResolver) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [path] [path2] [path3] [path...]",
		Short: "Resolve manifest files",
		Long: `Resolve manifest files. If a directory is inputted all manifest files without a lock file are resolved.
Example:
$ debricked resolve go.mod pkg/
`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(resolver),
	}
	fileExclusionExample := filepath.Join("*", "**.lock")
	dirExclusionExample := filepath.Join("**", "node_modules", "**")
	exampleFlags := fmt.Sprintf("-e \"%s\" -e \"%s\"", fileExclusionExample, dirExclusionExample)
	cmd.Flags().StringArrayVarP(&exclusions, ExclusionFlag, "e", exclusions, `The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Exclude flags could alternatively be set using DEBRICKED_EXCLUSIONS="path1,path2,path3".

Example: 
$ debricked resolve . `+exampleFlags)
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
	verboseDoc := strings.Join(
		[]string{
			"This flag allows you to reduce error output for resolution.",
			"\nExample:\n$ debricked resolve --verbose=false",
		}, "\n")
	cmd.Flags().BoolVar(&verbose, VerboseFlag, true, verboseDoc)
	npmPreferredDoc := strings.Join(
		[]string{
			"This flag allows you to select which package manager will be used as a resolver: Yarn (default) or NPM.",
			"Example: debricked resolve --prefer-npm",
		}, "\n")

	cmd.Flags().BoolP(NpmPreferredFlag, "", npmPreferred, npmPreferredDoc)

	cmd.Flags().IntVarP(&strictness, ResolutionStrictFlag, "s", file.StrictAll, `Allows you to configure exit code 1 or 0 depending on if the resolution was successful or not.
Strictness Level | Meaning
---------------- | -------
0 (default)      | Always exit with code 0, even if any or all files failed to resolve
1                | Exit with code 1 if all files failed to resolve, otherwise exit with code 0
2                | Exit with code 1 if any file failed to resolve, otherwise exit with code 0
3                | Exit with code 1 if all files failed to resolve, if any but not all files failed to resolve exit with code 3, otherwise exit with code 0
`)

	viper.MustBindEnv(ExclusionFlag)
	viper.MustBindEnv(NpmPreferredFlag)

	return cmd
}

func RunE(resolver resolution.IResolver) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = append(args, ".")
		}
		options := resolution.DebrickedOptions{
			Exclusions:   viper.GetStringSlice(ExclusionFlag),
			Verbose:      viper.GetBool(VerboseFlag),
			Regenerate:   viper.GetInt(RegenerateFlag),
			NpmPreferred: viper.GetBool(NpmPreferredFlag),
		}
		_, err := resolver.Resolve(args, options)

		return err
	}
}
