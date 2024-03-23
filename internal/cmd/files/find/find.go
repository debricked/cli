package find

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exclusions = file.Exclusions()
var inclusions []string
var jsonPrint bool
var lockfileOnly bool
var strictness int

const (
	ExclusionFlag    = "exclusion"
	InclusionFlag    = "inclusion"
	JsonFlag         = "json"
	LockfileOnlyFlag = "lockfile"
	StrictFlag       = "strict"
)

func NewFindCmd(finder file.IFinder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [path]",
		Short: "Find all dependency files in inputted path",
		Long: `Find all dependency files in inputted path. Related files are grouped together. 
For example ` + "`package.json`" + ` with ` + "`package-lock.json`.",
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

Exclude flags could alternatively be set using DEBRICKED_EXCLUSIONS="path1,path2,path3".

Example: 
$ debricked files find . `+exampleFlags)
	cmd.Flags().StringArrayVar(
		&inclusions,
		InclusionFlag,
		[]string{},
		`Forces inclusion of specified terms, see exclusion flag for more information on supported terms.
Examples: 
$ debricked scan . --include /node_modules/`)
	cmd.Flags().BoolVarP(&jsonPrint, JsonFlag, "j", false, `Print files in JSON format
Format:
[
  {
    "manifestFile": "package.json",
    "lockFiles": [
      "yarn.lock"
    ]
  },
]
`)
	cmd.Flags().BoolVarP(&lockfileOnly, LockfileOnlyFlag, "l", false, "If set, only lock files are found")
	cmd.Flags().IntVarP(&strictness, StrictFlag, "s", file.StrictAll, `Allows to control which files will be matched:
Strictness Level | Meaning
---------------- | -------
0 (default)      | Returns all matched manifest and lock files regardless if they're paired or not
1                | Returns only lock files and pairs of manifest and lock file
2                | Returns only pairs of manifest and lock file
`)

	viper.MustBindEnv(ExclusionFlag)
	viper.MustBindEnv(InclusionFlag)
	viper.MustBindEnv(JsonFlag)
	viper.MustBindEnv(LockfileOnlyFlag)
	viper.MustBindEnv(StrictFlag)

	return cmd
}

func RunE(f file.IFinder) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		err := AssertFlagsAreValid()
		if err != nil {
			return err
		}

		fileGroups, err := f.GetGroups(
			path,
			viper.GetStringSlice(ExclusionFlag),
			viper.GetStringSlice(InclusionFlag),
			viper.GetBool(LockfileOnlyFlag),
			viper.GetInt(StrictFlag),
		)
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

func AssertFlagsAreValid() error {
	if viper.GetBool(LockfileOnlyFlag) && viper.GetInt(StrictFlag) != file.StrictAll {
		return errors.New("'lockfile' and 'strict' flags are mutually exclusive")
	}

	if viper.GetInt(StrictFlag) < file.StrictAll || viper.GetInt(StrictFlag) > file.StrictPairs {
		return errors.New("'strict' supports values within range 0-2")
	}

	return nil
}
