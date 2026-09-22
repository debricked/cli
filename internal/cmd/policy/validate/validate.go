package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/policy"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var output string

const OutputFlag = "output"

const (
	OutputText = "text"
	OutputJson = "json"
)

const (
	// ValidationExitCode is returned when the policy file has validation errors
	ValidationExitCode = 1
	// FailureExitCode is returned when the validation could not be performed
	FailureExitCode = 2
)

var ErrInvalidPolicyFile = errors.New("policy file is invalid")

func NewValidateCmd(validator policy.IValidator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate a Debricked policy file",
		Long: `Validate a ` + policy.DefaultFileName + ` file against Debricked's policy schema.

If no path is given, ` + policy.DefaultFileName + ` in the current directory is validated.
If a directory is given, the ` + policy.DefaultFileName + ` inside it is validated.

Exit codes:
Code | Meaning
---- | -------
0    | The policy file is valid
1    | The policy file contains validation errors
2    | The validation could not be performed, for example due to a missing file or an unreachable service

Examples:
$ debricked policy validate
$ debricked policy validate .debricked/debricked_policy.json
$ debricked policy validate --output json
`,
		Args: cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(validator),
	}

	cmd.Flags().StringVarP(&output, OutputFlag, "o", OutputText, `Set output format.

Supported options are: 'text', 'json'`)
	viper.MustBindEnv(OutputFlag)

	return cmd
}

func RunE(v policy.IValidator) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		format := strings.ToLower(strings.TrimSpace(viper.GetString(OutputFlag)))
		if format != OutputText && format != OutputJson {
			return fail(cmd, fmt.Errorf("unsupported output format \"%s\". Supported options are '%s' and '%s'", format, OutputText, OutputJson))
		}

		var path string
		if len(args) > 0 {
			path = args[0]
		}

		result, err := v.Validate(path)
		if err != nil {
			return fail(cmd, err)
		}

		if err = print(cmd.OutOrStdout(), result, format); err != nil {
			return fail(cmd, err)
		}

		if !result.Valid {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			return cmderror.CommandError{Code: ValidationExitCode, Err: ErrInvalidPolicyFile}
		}

		return nil
	}
}

func print(writer io.Writer, result policy.Result, format string) error {
	if format == OutputJson {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")

		return encoder.Encode(result)
	}

	printText(writer, result)

	return nil
}

func printText(writer io.Writer, result policy.Result) {
	if result.Valid {
		_, _ = fmt.Fprintf(writer, "%s %s is valid\n", color.GreenString("✔"), result.File)

		return
	}

	errorLabel := "errors"
	if len(result.Errors) == 1 {
		errorLabel = "error"
	}
	_, _ = fmt.Fprintf(
		writer,
		"%s %s has %d validation %s\n",
		color.RedString("⨯"),
		result.File,
		len(result.Errors),
		errorLabel,
	)
	for _, validationError := range result.Errors {
		_, _ = fmt.Fprintf(writer, "%s\n", validationError.String())
	}
}

// fail reports a non-validation failure. The message goes to stderr to keep
// stdout consumable by CI integrations.
func fail(cmd *cobra.Command, err error) error {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s %s\n", color.RedString("⨯"), err.Error())

	return cmderror.CommandError{Code: FailureExitCode, Err: err}
}
