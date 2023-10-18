package resolve

import (
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/resolution"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exclusions = file.Exclusions()

const (
	ExclusionFlag = "exclusion"
)

func NewResolveCmd(resolver resolution.IResolver) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [path]",
		Short: "Resolve manifest files",
		Long: `Resolve manifest files. If a directory is inputted all manifest files without a lock file are resolved.
Example:
$ debricked files resolve go.mod pkg/
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
$ debricked files resolve . `+exampleFlags)

	viper.MustBindEnv(ExclusionFlag)

	return cmd
}

func RunE(resolver resolution.IResolver) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = append(args, ".")
		}
		_, err := resolver.Resolve(args, viper.GetStringSlice(ExclusionFlag))

		return err
	}
}
