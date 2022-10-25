package find

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var exclusions []string
var jsonPrint bool
var lockfileOnly bool

const (
	ExclusionsFlag   = "exclusions"
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
		RunE: RunE(finder),
	}
	cmd.Flags().StringArrayVarP(&exclusions, ExclusionsFlag, "e", exclusions, `The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Examples: 
$ debricked files find . -e "*/**.lock" -e "**/node_modules/**" 
$ debricked files find . -e "*\**.exe" -e "**\node_modules\**" 
`)
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
	_ = viper.BindPFlags(cmd.Flags())
	viper.MustBindEnv(ExclusionsFlag)
	viper.MustBindEnv(JsonFlag)
	viper.MustBindEnv(LockfileOnlyFlag)

	return cmd
}

func RunE(f file.IFinder) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		directoryPath := args[0]
		fileGroups, err := f.GetGroups(directoryPath, viper.GetStringSlice(ExclusionsFlag), viper.GetBool(LockfileOnlyFlag))
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
