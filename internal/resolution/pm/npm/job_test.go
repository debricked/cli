package npm

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/npm/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunInstallCmdErr(t *testing.T) {
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
			error: "npm ERR! code ETARGET\nnpm ERR! notarget No matching version found for chalk@^113.0.0.\nnpm ERR! notarget In most cases you or one of your dependencies are requesting\nnpm ERR! notarget a package version that doesn't exist.",
			doc:   "Failed to find package \"chalk@^113.0.0\" that satisfies the requirement from package.json file. In most cases you or one of your dependencies are requesting a package version that doesn't exist. Please check that package versions are correct in your package.json file.",
		},
		{
			name:  "Invalid package name or private package",
			error: "npm ERR! code E404\nnpm ERR! 404 Not Found - GET https://registry.npmjs.org/chalke - Not found\nnpm ERR! 404 \nnpm ERR! 404  'chalke@^3.0.0' is not in this registry.\nnpm ERR! 404 You should bug the author to publish it (or use the name yourself!)\nnpm ERR! 404 \nnpm ERR! 404 Note that you can also install from a\nnpm ERR! 404 tarball, folder, http url, or git url.",
			doc:   "Failed to find package \"chalke@^3.0.0\" that satisfies the requirement from dependencies. Please check that dependencies are correct in your package.json file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "No internet connection",
			error: "npm ERR! code EAI_AGAIN\nnpm ERR! syscall getaddrinfo\nnpm ERR! errno EAI_AGAIN\nnpm ERR! request to https://registry.npmjs.org/chalke failed, reason: getaddrinfo EAI_AGAIN registry.npmjs.org",
			doc:   "Package registry \"registry.npmjs.org\" is not available at the moment. There might be a trouble with your network connection.",
		},
		{
			name:  "Permission denied",
			error: "npm ERR! Error: EACCES, open '/home/me/.npm/semver/3.0.1/package/package.json'\nnpm ERR!  { [Error: EACCES, open '/home/me/.npm/semver/3.0.1/package/package.json']\nnpm ERR!   errno: 3,\nnpm ERR!   code: 'EACCES',\nnpm ERR!   path: '/home/me/.npm/semver/3.0.1/package/package.json',\nnpm ERR!   parent: 'gulp' }\nnpm ERR! \nnpm ERR! Please try running this command again as root/Administrator.\n",
			doc:   "Couldn't get access to \"/home/me/.npm/semver/3.0.1/package/package.json\". Please check permissions or try running this command again as root/Administrator.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmdErr := errors.New(c.error)
			cmdFactoryMock := testdata.NewEchoCmdFactory()
			cmdFactoryMock.MakeInstallErr = cmdErr
			cmd, _ := cmdFactoryMock.MakeInstallCmd("echo", "package.json")

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)
			expectedError.SetStatus("installing dependencies")
			expectedError.SetCommand(cmd.String())

			j := NewJob("file", true, cmdFactoryMock)

			go jobTestdata.WaitStatus(j)
			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, j.Errors().GetAll(), 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	j := NewJob("file", true, cmdMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}
