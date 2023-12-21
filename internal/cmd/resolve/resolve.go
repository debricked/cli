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
)

const (
	ExclusionFlag    = "exclusion"
	VerboseFlag      = "verbose"
	NpmPreferredFlag = "prefer-npm"
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
	cmd.Flags().BoolVar(&verbose, VerboseFlag, true, "set to false to disable extensive resolution error messages")

	npmPreferredDoc := strings.Join(
		[]string{
			"This flag allows you to select which package manager will be used as a resolver: Yarn (default) or NPM.",
			"Example: debricked resolve --prefer-npm",
		}, "\n")

	cmd.Flags().BoolP(NpmPreferredFlag, "", npmPreferred, npmPreferredDoc)

	viper.MustBindEnv(ExclusionFlag)
	viper.MustBindEnv(NpmPreferredFlag)

	return cmd
}

func RunE(resolver resolution.IResolver) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = append(args, ".")
		}

		resolver.SetNpmPreferred(viper.GetBool(NpmPreferredFlag))
		_, err := resolver.Resolve(args, viper.GetStringSlice(ExclusionFlag), viper.GetBool(VerboseFlag))

		return err
	}
}
