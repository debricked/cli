package find

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var exclusions = file.DefaultExclusions()
var jsonPrint bool
var lockfileOnly bool

const (
	ExclusionFlag    = "exclusion"
	JsonFlag         = "json"
	LockfileOnlyFlag = "lockfile"
)

func NewFindCmd(finder file.IFinder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [path]",
		Short: "Find all dependency files in inputted path",
		Long: `Find all dependency files in inputted path. Related files are grouped together. 
For example ` + "`package.json`" + ` with ` + "`package-lock.json`.",
		Args: validateArgs,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(finder),
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

Example: 
$ debricked files find . `+exampleFlags)

	cmd.Flags().BoolVarP(&jsonPrint, JsonFlag, "j", false, `Print files in JSON format
Format:
[
  {
    "dependencyFile": "package.json",
    "lockFiles": [
      "yarn.lock"
    ]
  },
]
`)
	cmd.Flags().BoolVarP(&lockfileOnly, LockfileOnlyFlag, "l", false, "If set, only lock files are found")

	viper.MustBindEnv(ExclusionFlag)
	viper.MustBindEnv(JsonFlag)
	viper.MustBindEnv(LockfileOnlyFlag)

	return cmd
}

func RunE(f file.IFinder) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		directoryPath := args[0]
		fileGroups, err := f.GetGroups(directoryPath, viper.GetStringSlice(ExclusionFlag), viper.GetBool(LockfileOnlyFlag))
		if err != nil {
			return err
		}
		if viper.GetBool(JsonFlag) {
			jsonFileGroups, _ := json.Marshal(fileGroups.ToSlice())
			fmt.Println(string(jsonFileGroups))
		} else {
			for _, group := range fileGroups.ToSlice() {
				group.Print()
			}
		}

		return nil
	}
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires path")
	}
	if isValidFilepath(args[0]) {
		return nil
	}
	return fmt.Errorf("invalid path specified: %s", args[0])
}

func isValidFilepath(path string) bool {
	_, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	return true
}
