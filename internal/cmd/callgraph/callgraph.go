package callgraph

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/callgraph"
	cg "github.com/debricked/cli/internal/callgraph"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	exclusions         = file.DefaultExclusions()
	inclusions         []string
	buildDisabled      bool
	generateTimeout    int
	languages          string
	supportedLanguages = []string{"java", "golang"}
	languageMap        = map[string]string{
		"java":   "maven",
		"golang": "go",
	}
)

const (
	ExclusionFlag       = "exclusion"
	InclusionFlag       = "inclusion"
	NoBuildFlag         = "no-build"
	GenerateTimeoutFlag = "generate-timeout"
	LanguagesFlag       = "languages"
)

func NewCallgraphCmd(generator cg.IGenerator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "callgraph [path]",
		Short: "Generate a static call graph for the given directory and subdirectories",
		Long: `Generate a static call graph for a project in the given directory. The command consists of two main parts: build and callgraph. 
Build: Build the project and resolve dependencies. In this step, all necessary .class files are created.
Callgraph: Generate the static call graph using debricked Reachability Analysis.

The full documentation is available here https://docs.debricked.com/tools-and-integrations/cli/debricked-cli

Example:
$ debricked callgraph 
`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(generator),
	}
	fileExclusionExample := filepath.Join("'*", "**.lock'")
	dirExclusionExample := filepath.Join("'**", "node_modules", "**'")
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
$ debricked scan . --include '**/node_modules/**'`)
	cmd.Flags().BoolVar(&buildDisabled, NoBuildFlag, false, `Do not automatically build all source code in the project to enable call graph generation.
This option requires a pre-built project. For more detailed documentation on the callgraph generation, visit:
https://docs.debricked.com/tools-and-integrations/cli/debricked-cli#callgraph`)
	cmd.Flags().IntVar(&generateTimeout, GenerateTimeoutFlag, 60*60, "Timeout (in seconds) on call graph generation.")
	cmd.Flags().StringVarP(&languages, LanguagesFlag, "l", strings.Join(supportedLanguages, ","), "Colon separated list of languages to create a call graph for.")

	viper.MustBindEnv(ExclusionFlag)

	return cmd
}

func parseAndValidateLanguages(languages string) ([]string, error) {
	if languages == "" {
		// default to all supported languages
		return supportedLanguages, nil
	}

	parsedLanguages := strings.Split(languages, ",")

	for _, language := range parsedLanguages {
		if !isSupportedLanguage(language) {
			return nil, errors.New(language + " is not a supported language")
		}
	}

	return parsedLanguages, nil
}

func isSupportedLanguage(language string) bool {
	for _, supportedLanguage := range supportedLanguages {
		if language == supportedLanguage {
			return true
		}
	}

	return false
}

func RunE(callgraph callgraph.IGenerator) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = append(args, ".")
		}

		languages, err := parseAndValidateLanguages(languages)
		if err != nil {
			return err
		}

		configs := []conf.IConfig{}
		version := viper.GetString("cliVersion")

		for _, language := range languages {
			configs = append(configs, conf.NewConfig(language, args, map[string]string{}, !buildDisabled, languageMap[language], version))
		}

		options := cg.DebrickedOptions{
			Paths:      args,
			Exclusions: viper.GetStringSlice(ExclusionFlag),
			Inclusions: viper.GetStringSlice(InclusionFlag),
			Configs:    configs,
			Timeout:    viper.GetInt(GenerateTimeoutFlag),
		}

		return callgraph.GenerateWithTimer(options)
	}
}
