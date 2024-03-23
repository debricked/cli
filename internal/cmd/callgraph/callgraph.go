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

var (
	exclusions      = file.DefaultExclusions()
	inclusions      []string
	buildDisabled   bool
	generateTimeout int
)

const (
	ExclusionFlag       = "exclusion"
	InclusionFlag       = "inclusion"
	NoBuildFlag         = "no-build"
	GenerateTimeoutFlag = "generate-timeout"
)

func NewCallgraphCmd(generator callgraph.IGenerator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "callgraph [path]",
		Short: "Generate a static call graph for the given directory and subdirectories",
		Long: `Generate a static call graph for a project in the given directory. The command consists of two main parts: build and callgraph. 
Build: Build the project and resolve dependencies. In this step, all necessary .class files are created.
Callgraph: Generate the static call graph using debricked Reachability Analysis.

The full documentation is available here https://portal.debricked.com/debricked-cli-63/debricked-cli-documentation-298

Example:
$ debricked callgraph 
`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(generator),
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
$ debricked callgraph . `+exampleFlags)
	cmd.Flags().StringArrayVar(
		&inclusions,
		InclusionFlag,
		[]string{},
		`Forces inclusion of specified terms, see exclusion flag for more information on supported terms.
Examples: 
$ debricked scan . --include /node_modules/`)
	cmd.Flags().BoolVar(&buildDisabled, NoBuildFlag, false, `Do not automatically build all source code in the project to enable call graph generation.
This option requires a pre-built project. For more detailed documentation on the callgraph generation, visit our portal:
https://portal.debricked.com/debricked-cli-63/debricked-cli-documentation-298?tid=298&fid=63#callgraph`)
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

		err := callgraph.GenerateWithTimer(args, viper.GetStringSlice(ExclusionFlag), viper.GetStringSlice(InclusionFlag), configs, viper.GetInt(GenerateTimeoutFlag))

		return err
	}
}
