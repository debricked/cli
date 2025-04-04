package sbt

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/sbt/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{}, BuildService{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErr(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "General error",
			error: "cmd-error",
			doc:   util.UnknownError,
		},
		{
			name:  "SBT not found",
			error: "        |exec: \"sbt\": executable file not found in $PATH",
			doc:   "SBT wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "POM generation error",
			error: " |[error] Error occurred while processing command: makePom",
			doc:   "Failed to generate Maven POM file. SBT encountered an error during the makePom task. Error details: Error occurred while generating the POM file",
		},
		{
			name:  "Build file not found",
			error: " |[error] not found: /home/user/project/build.sbt",
			doc:   "SBT configuration file not found. Please ensure that your project contains a valid build.sbt file. Error details: build.sbt file not found",
		},
		{
			name:  "Invalid build file",
			error: " |[error] Illegal character in build file at line 15",
			doc:   "Failed to parse SBT build file. Your build.sbt file contains syntax errors. Please check the build file for errors",
		},
		{
			name:  "No Internet",
			error: " |[error] Connection timed out: connect",
			doc:   "We weren't able to retrieve one or more dependencies or plugins. Please check your Internet connection and try again.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)

			cmdErr := errors.New(c.error)
			j := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr}, testdata.BuildServiceMock{})

			go jobTestdata.WaitStatus(j)

			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, allErrors, 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunCmdOutputErr(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"}, testdata.BuildServiceMock{})

	go jobTestdata.WaitStatus(j)

	j.Run()

	error := j.Errors()
	assert.True(t, error.HasError())
	
	allErrors := error.GetAll()
	assert.Len(t, allErrors,2)
	assert.Contains(t, allErrors[0].Error(), "executable file not found")
	assert.Contains(t, allErrors[0].Documentation(), "SBT wasn't found")
}

func TestRunCmdOutputErrNoOutput(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "go", Arg: "bad-arg"}, testdata.BuildServiceMock{})

	go jobTestdata.WaitStatus(j)

	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 2)
	err := errs[0]

	assert.Contains(t, err.Error(), "unknown command")
}

func TestRun(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"}, testdata.BuildServiceMock{})

	go jobTestdata.WaitStatus(j)

	j.Run()
	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.Equal(t, 1, len(errs))

	assert.True(t, j.Errors().HasError())
}

func TestRunWithBuildServiceError(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "empty file",
			error: "EOF",
			doc:   "This file doesn't contain valid SBT build content",
		},
		{
			name:  "syntax error",
			error: "syntax error in build.sbt",
			doc:   "syntax error in build.sbt",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"}, testdata.BuildServiceMock{
				Err: errors.New(c.error),
			})

			go jobTestdata.WaitStatus(j)

			j.Run()

			allErrors := j.Errors().GetAll()

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetStatus("parsing SBT build file")
			expectedError.SetDocumentation(c.doc)

			assert.Len(t, allErrors, 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}
