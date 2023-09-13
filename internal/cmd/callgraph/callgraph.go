package callgraph

import (
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exclusions = file.DefaultExclusions()

const (
	ExclusionFlag       = "exclusion"
	NoBuildFlag         = "no-build"
	GenerateTimeoutFlag = "generate-timeout"
)

var buildDisabled bool
var generateTimeout int

func NewCallgraphCmd(generator callgraph.IGenerator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "callgraph [path]",
		Short: "Generate a static call graph for the given directory and subdirectories",
		Long: `Generate a static call graph for a project. If a directory is inputted all manifest files without a lock file are resolved. 

Example:
$ debricked callgraph 
`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(generator),
	}
	// TODO: add to docs: Complete documentation with advanced user guide is available at https://portal.debricked.com/docs/...
	fileExclusionExample := filepath.Join("*", "**.lock")
	dirExclusionExample := filepath.Join("**", "node_modules", "**")
	exampleFlags := fmt.Sprintf("-e \"%s\" -e \"%s\"", fileExclusionExample, dirExclusionExample)

	cmd.Flags().StringArrayVarP(&exclusions, ExclusionFlag, "e", exclusions, 
`The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Example: 
$ debricked files resolve . `+exampleFlags)
	cmd.Flags().BoolVar(&buildDisabled, NoBuildFlag, false, "Do not automatically build all source code in the project to enable call graph generation. This option requires a pre-built project.")
	cmd.Flags().IntVar(&generateTimeout, GenerateTimeoutFlag, 60*60, "Timeout (in seconds) on call graph generation.")

	viper.MustBindEnv(ExclusionFlag)

	return cmd
}

func RunE(callgraph callgraph.IGenerator) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = append(args, ".")
		}
		configs := []conf.IConfig{
			conf.NewConfig("java", []string{}, map[string]string{}, !viper.GetBool(NoBuildFlag), "maven"),
		}

		err := callgraph.GenerateWithTimer(args, viper.GetStringSlice(ExclusionFlag), configs, viper.GetInt(GenerateTimeoutFlag))
		// err := callgraph.Generate(args, viper.GetStringSlice(ExclusionFlag), configs, nil)

		return err
	}
}