package yarn

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/yarn/testdata"
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
			name:  "Yarn not found",
			error: "        |exec: \"yarn\": executable file not found in $PATH",
			doc:   "Yarn wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Invalid package name",
			error: "error Error: https://registry.yarnpkg.com/chalke: Not found\n    at Request.params.callback [as _callback] (/usr/local/lib/node_modules/yarn/lib/cli.js:66148:18)",
			doc:   "Failed to find package \"https://registry.yarnpkg.com/chalke\" that satisfies the requirement from yarn dependencies. Please check that dependencies are correct in your package.json file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "Invalid package name (legacy)",
			error: `error An unexpected error occurred: "https://registry.yarnpkg.com/chalke: Not found".`,
			doc:   "Failed to find package \"https://registry.yarnpkg.com/chalke\" that satisfies the requirement from yarn dependencies. Please check that dependencies are correct in your package.json file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "Invalid package version",
			error: "error Couldn't find any versions for \"chalk\" that matches \"^300.0.0\"\ninfo Visit https://yarnpkg.com/en/docs/cli/install for documentation about this command.",
			doc:   "Couldn't find any versions for \"chalk\" that matches \"^300.0.0\". Please check that dependencies are correct in your package.json file.",
		},
		{
			name:  "No internet connection",
			error: "error Error: getaddrinfo ENOTFOUND nexus.dev\n    at GetAddrInfoReqWrap.onlookup [as oncomplete] (dns.js:66:26)\n",
			doc:   "Package registry \"nexus.dev\" is not available at the moment. There might be a trouble with your network connection.",
		},
		{
			name:  "Private package",
			error: "Error: https://registry.npmjs.org/@private/my-private-package/-/my-private-package-0.0.5.tgz: Request failed \"404 Not Found\"",
			doc:   "Failed to find a package that satisfies requirements for yarn dependencies: https://registry.npmjs.org/@private/my-private-package/-/my-private-package-0.0.5.tgz. This could mean that the package or version does not exist or is private.\n If this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
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
