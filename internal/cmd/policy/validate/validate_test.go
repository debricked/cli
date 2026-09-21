package validate

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/cmd/policy/testdata"
	"github.com/debricked/cli/internal/policy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const policyFile = "debricked_policy.json"

func setUp(t *testing.T, validator policy.IValidator, format string) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	cmd := NewValidateCmd(validator)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	viper.Set(OutputFlag, format)
	t.Cleanup(func() {
		viper.Set(OutputFlag, "")
	})

	return cmd, stdout, stderr
}

func run(cmd *cobra.Command, validator policy.IValidator, args ...string) error {
	return RunE(validator)(cmd, args)
}

func TestNewValidateCmd(t *testing.T) {
	cmd := NewValidateCmd(&testdata.ValidatorMock{})

	assert.Empty(t, cmd.Commands())
	flag := cmd.Flags().Lookup(OutputFlag)
	assert.NotNil(t, flag)
	assert.Equal(t, "o", flag.Shorthand)
	assert.Equal(t, OutputText, flag.DefValue)

	match := false
	for _, key := range viper.AllKeys() {
		if key == OutputFlag {
			match = true

			break
		}
	}
	assert.Truef(t, match, "failed to assert that %s was present", OutputFlag)
}

func TestValidateCmdArgs(t *testing.T) {
	cmd := NewValidateCmd(&testdata.ValidatorMock{})

	assert.NoError(t, cmd.Args(cmd, []string{}))
	assert.NoError(t, cmd.Args(cmd, []string{"debricked_policy.json"}))
	assert.Error(t, cmd.Args(cmd, []string{"a.json", "b.json"}))
}

func TestRunEValidFile(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{File: policyFile, Valid: true, Errors: []policy.ValidationError{}},
	}
	cmd, stdout, stderr := setUp(t, validator, OutputText)

	err := run(cmd, validator)

	assert.NoError(t, err)
	assert.Contains(t, stdout.String(), policyFile+" is valid")
	assert.Empty(t, stderr.String())
	assert.Equal(t, "", validator.Path)
}

func TestRunEPathArgument(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{File: "policies/debricked_policy.json", Valid: true},
	}
	cmd, _, _ := setUp(t, validator, OutputText)

	err := run(cmd, validator, "policies")

	assert.NoError(t, err)
	assert.Equal(t, "policies", validator.Path)
}

func TestRunEValidationErrors(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{
			File: policyFile,
			Errors: []policy.ValidationError{
				{PropertyPath: "policies[0].name", Message: "The property name is required"},
				{Message: "Array must have at least 1 item"},
			},
		},
	}
	cmd, stdout, stderr := setUp(t, validator, OutputText)

	err := run(cmd, validator)

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, ValidationExitCode, cmdErr.Code)
	assert.ErrorIs(t, err, ErrInvalidPolicyFile)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.Empty(t, stderr.String())

	output := stdout.String()
	assert.Contains(t, output, policyFile+" has 2 validation errors")
	assert.Contains(t, output, "policies[0].name: The property name is required\n")
	assert.Contains(t, output, "Array must have at least 1 item\n")
}

func TestRunESingleValidationError(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{
			File:   policyFile,
			Errors: []policy.ValidationError{{Message: "Array must have at least 1 item"}},
		},
	}
	cmd, stdout, _ := setUp(t, validator, OutputText)

	err := run(cmd, validator)

	assert.Error(t, err)
	assert.Contains(t, stdout.String(), "has 1 validation error\n")
}

func TestRunEJsonOutput(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{
			File: policyFile,
			Errors: []policy.ValidationError{
				{PropertyPath: "policies[0].name", Message: "The property name is required"},
			},
		},
	}
	cmd, stdout, _ := setUp(t, validator, OutputJson)

	err := run(cmd, validator)

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, ValidationExitCode, cmdErr.Code)

	var result policy.Result
	assert.NoError(t, json.Unmarshal(stdout.Bytes(), &result))
	assert.False(t, result.Valid)
	assert.Equal(t, policyFile, result.File)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "policies[0].name", result.Errors[0].PropertyPath)
	assert.Equal(t, "The property name is required", result.Errors[0].Message)
}

func TestRunEJsonOutputValidFile(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{File: policyFile, Valid: true, Errors: []policy.ValidationError{}},
	}
	cmd, stdout, _ := setUp(t, validator, "JSON ")

	err := run(cmd, validator)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"file": "debricked_policy.json", "valid": true, "errors": []}`, stdout.String())
}

func TestRunEValidationFailure(t *testing.T) {
	validator := &testdata.ValidatorMock{Error: policy.ErrNoPolicyFile}
	cmd, stdout, stderr := setUp(t, validator, OutputText)

	err := run(cmd, validator)

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, FailureExitCode, cmdErr.Code)
	assert.ErrorIs(t, err, policy.ErrNoPolicyFile)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.Empty(t, stdout.String(), "failed to assert that stdout was kept clean for CI integrations")
	assert.Contains(t, stderr.String(), policy.ErrNoPolicyFile.Error())
}

func TestRunEUnsupportedOutputFormat(t *testing.T) {
	validator := &testdata.ValidatorMock{}
	cmd, _, stderr := setUp(t, validator, "yaml")

	err := run(cmd, validator)

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, FailureExitCode, cmdErr.Code)
	assert.Contains(t, stderr.String(), "unsupported output format")
	assert.Equal(t, "", validator.Path, "failed to assert that validation was skipped")
}

func TestRunEPrintError(t *testing.T) {
	validator := &testdata.ValidatorMock{
		Result: policy.Result{
			File:   policyFile,
			Valid:  true,
			Errors: []policy.ValidationError{{Context: map[string]interface{}{"invalid": make(chan int)}}},
		},
	}
	cmd, _, stderr := setUp(t, validator, OutputJson)

	err := run(cmd, validator)

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, FailureExitCode, cmdErr.Code)
	assert.Contains(t, stderr.String(), "json")
}

func TestPreRun(t *testing.T) {
	cmd := NewValidateCmd(&testdata.ValidatorMock{})
	cmd.PreRun(cmd, nil)
}

func TestRunEErrorIsNotWrappedTwice(t *testing.T) {
	validationErr := errors.New("boom")
	validator := &testdata.ValidatorMock{Error: validationErr}
	cmd, _, _ := setUp(t, validator, OutputText)

	err := run(cmd, validator)

	assert.ErrorIs(t, err, validationErr)
}
