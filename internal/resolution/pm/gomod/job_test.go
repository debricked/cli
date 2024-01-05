package gomod

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/gomod/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/internal/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunGraphCmdErr(t *testing.T) {
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
			name:  "Invalid package version",
			error: "        |go: errors parsing go.mod:\n        |/home/asus/Projects/playground/imessage/go.mod:8:2: require github.com/google/uuid: version \"v111.5.0\" invalid: should be v0 or v1, not v111",
			doc:   "Failed to find package \"github.com/google/uuid\". Please check that package versions are correct in the manifest file. It should be v0 or v1, not v111.",
		},
		{
			name:  "Invalid package name",
			error: "        |go: github.com/google/yuid@v1.5.0: reading github.com/google/yuid/go.mod at revision v1.5.0: git ls-remote -q origin in /home/asus/go/pkg/mod/cache/vcs/0292faa5faa65b4148fe687f4ad2478601180035651bd75864518f9b0a6ddd2c: exit status 128:\n        |       fatal: could not read Username for 'https://github.com': terminal prompts disabled\n        |Confirm the import path was entered correctly.",
			doc:   "Failed to find package \"github.com/google/yuid@v1.5.0\" that satisfies the requirements. Please check that dependencies are correct in the manifest file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "No package in repository",
			error: "github.com/MyCompany/mypackage: module github.com/MyCompany/mypackage@latest found (v0.0.0-20220126203606-a88fea44b771), but does not contain package github.com/MyCompany/mypackage",
			doc:   "We weren't able to find a package in provided repository \"github.com/MyCompany/mypackage\". Please check that repository address is spelled correct and it actually contains a Go package.",
		},
		{
			name:  "Unable to resolve package",
			error: "go: module github.com/Private-Org/go-commons: git ls-remote -q origin in /Users/arieroos/Workspace/go/pkg/mod/cache/vcs/4371bddb1b5a61a8f85ece4f86eaa40bbac4cc02925be418880bcce25aafa433: exit status 128:\n        git@github.com: Permission denied (publickey).\n        fatal: Could not read from remote repository.",
			doc:   "Failed to find package \"github.com/Private-Org/go-commons\" that satisfies the requirements. Please check that dependencies are correct in the manifest file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "No internet connection",
			error: "        |go: github.com/google/yuid@v1.5.0: Get \"https://proxy.golang.org/github.com/google/yuid/@v/v1.5.0.mod\": dial tcp: lookup proxy.golang.org on 127.0.0.53:53: server misbehaving",
			doc:   "Registry \"proxy.golang.org\" is not available at the moment. There might be a trouble with your network connection.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmdErr := errors.New(c.error)
			cmdFactoryMock := testdata.NewEchoCmdFactory()
			cmdFactoryMock.MakeGraphCmdErr = cmdErr
			cmd, _ := cmdFactoryMock.MakeGraphCmd("echo")

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)
			expectedError.SetStatus("creating dependency graph")
			expectedError.SetCommand(cmd.String())

			j := NewJob("file", cmdFactoryMock, nil)

			go jobTestdata.WaitStatus(j)
			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, j.Errors().GetAll(), 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.GraphCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListCmdErr = cmdErr
	j := NewJob("file", cmdFactoryMock, nil)

	cmd, _ := cmdFactoryMock.MakeListCmd("echo")
	expectedError := util.NewPMJobError(cmdErr.Error())
	expectedError.SetStatus("creating dependency version list")
	expectedError.SetCommand(cmd.String())

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, allErrors, expectedError)
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.ListCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	cmd, _ := cmdFactoryMock.MakeListCmd("echo")
	expectedError := util.NewPMJobError(createErr.Error())
	expectedError.SetStatus("creating lock file")
	expectedError.SetCommand(cmd.String())

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, allErrors, expectedError)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	cmd, _ := cmdFactoryMock.MakeListCmd("echo")
	expectedError := util.NewPMJobError(writeErr.Error())
	expectedError.SetStatus("creating lock file")
	expectedError.SetCommand(cmd.String())

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, allErrors, expectedError)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(closeErr.Error()))
}

func TestRun(t *testing.T) {
	fileContents := []byte("MakeGraphCmd\n\nMakeListCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Empty(t, j.Errors().GetAll())
	assert.Equal(t, fileContents, fileWriterMock.Contents)
}
