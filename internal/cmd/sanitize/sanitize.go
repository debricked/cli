package sanitize

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/policy"
	"github.com/debricked/cli/internal/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ExclusionFlag = "exclusion"
	InclusionFlag = "inclusion"
	JsonFlag      = "json"

	// FailureExitCode is returned when at least one file could not be sanitized.
	FailureExitCode = 1
)

var (
	exclusions = file.Exclusions()
	inclusions []string
	jsonPrint  bool
)

func NewSanitizeCmd(sanitizer sanitize.ISanitizer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sanitize [path]",
		Short: "Write redacted copies of dependency files",
		Long: `Write redacted copies of the files found by ` + "`debricked files find`" + `, keeping only what the scan uses.

Copies are written to .debricked/` + sanitize.OutputDirName + ` in the given path, with the same file names and
relative paths. That folder is cleared on every run. Original files are never changed.

The sanitized files are printed in the same format as ` + "`debricked files find`" + `, with the paths of
the copies instead of the originals. Skipped and failed files, and a summary, go to stderr.

A map from each placeholder to the real name and path is written to .debricked/` + sanitize.MapFileName + `.
It stays on this machine: it is outside the folder of copies and must not be uploaded.

File types that sanitize doesn't handle are skipped and reported, and no copy of them is written.
` + policy.DefaultFileName + ` files are copied unchanged, after checking that they are valid JSON, and listed
with the other files.

Exit codes:
Code | Meaning
---- | -------
0    | All supported files were sanitized
1    | At least one file could not be sanitized, or the run itself failed`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(sanitizer),
	}

	cmd.Flags().StringArrayVarP(&exclusions, ExclusionFlag, "e", exclusions, `Paths to exclude, in the same format as "debricked files find".
Could alternatively be set using DEBRICKED_EXCLUSIONS="path1,path2,path3".`)
	cmd.Flags().StringArrayVar(&inclusions, InclusionFlag, []string{}, `Forces inclusion of specified terms, see the exclusion flag.`)
	cmd.Flags().BoolVarP(&jsonPrint, JsonFlag, "j", false, `Print sanitized files in JSON format, as "debricked files find" does
Format:
[
  {
    "manifestFile": ".debricked/sanitized/package.json",
    "lockFiles": [
      ".debricked/sanitized/yarn.lock"
    ]
  },
]
`)

	viper.MustBindEnv(ExclusionFlag)
	viper.MustBindEnv(InclusionFlag)
	viper.MustBindEnv(JsonFlag)

	return cmd
}

func RunE(sanitizer sanitize.ISanitizer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		return Run(cmd, sanitizer, sanitize.Options{
			RootPath:   path,
			Exclusions: viper.GetStringSlice(ExclusionFlag),
			Inclusions: viper.GetStringSlice(InclusionFlag),
		}, viper.GetBool(JsonFlag))
	}
}

// Run sanitizes, prints the sanitized files in the format of `debricked files find` and returns
// an error with FailureExitCode if any file failed. `debricked files find --sanitize` uses it too.
func Run(cmd *cobra.Command, sanitizer sanitize.ISanitizer, options sanitize.Options, asJson bool) error {
	result, err := sanitizer.Sanitize(options)
	if err != nil {
		cmd.SilenceUsage = true

		return cmderror.CommandError{Code: FailureExitCode, Err: err}
	}

	if err = printGroups(cmd.OutOrStdout(), result.Groups, asJson); err != nil {
		return cmderror.CommandError{Code: FailureExitCode, Err: err}
	}
	printSummary(cmd.ErrOrStderr(), result)

	if len(result.Failed) > 0 {
		cmd.SilenceUsage = true

		return cmderror.CommandError{
			Code: FailureExitCode,
			Err:  errors.New("some files could not be sanitized, so the output is incomplete"),
		}
	}

	return nil
}

// printGroups uses the output format of `debricked files find`.
func printGroups(out io.Writer, groups []file.Group, asJson bool) error {
	if asJson {
		encoded, err := json.Marshal(groups)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, string(encoded))

		return err
	}

	for _, group := range groups {
		prefix := ""
		if group.HasFile() {
			_, _ = fmt.Fprintln(out, group.ManifestFile)
			prefix = " * "
		}
		for _, lockFile := range group.LockFiles {
			_, _ = fmt.Fprintln(out, prefix+lockFile)
		}
	}

	return nil
}

// printSummary lists paths and fixed reasons only, so no removed value can appear in the output.
func printSummary(out io.Writer, result sanitize.Result) {
	printIssues(out, "Skipped", result.Skipped)
	printIssues(out, "Failed", result.Failed)

	_, _ = fmt.Fprintf(
		out,
		"Sanitized: %d, skipped: %d, failed: %d\n",
		len(result.Sanitized),
		len(result.Skipped),
		len(result.Failed),
	)
	_, _ = fmt.Fprintf(out, "Sanitized copies: %s\n", result.OutputDir)
	_, _ = fmt.Fprintf(out, "Local map, not for upload: %s\n", result.MapPath)
}

func printIssues(out io.Writer, title string, issues []sanitize.FileIssue) {
	if len(issues) == 0 {
		return
	}
	_, _ = fmt.Fprintf(out, "%s:\n", title)
	for _, issue := range issues {
		_, _ = fmt.Fprintf(out, "  %s: %s\n", issue.Path, issue.Reason)
	}
}
