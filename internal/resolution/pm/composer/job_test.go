package composer

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/composer/testdata"
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

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", false, cmdFactoryMock)

	_, err := j.runInstallCmd()
	assert.NoError(t, err)

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
		error string
		doc   string
	}{
		{
			error: "cmd-error",
			doc:   util.UnknownError,
		},
		{
			error: "\n\n PHP's phar extension is missing. Composer requires it to run. Enable the extension or recompile php without --disable-phar then try again.",
			doc:   "Failed to build Composer dependency tree. Your runtime environment is missing one or more Composer requirements. Check error message below for more details:\n\n \n\n PHP's phar extension is missing. Composer requires it to run. Enable the extension or recompile php without --disable-phar then try again.",
		},
		{
			error: "require.debricked is invalid, it should have a vendor name, a forward slash, and a package name",
			doc:   "Couldn't resolve dependency debricked , please make sure it is spelt correctly:\n",
		},
		{
			error: "The following exception probably indicates you have misconfigured DNS resolver(s)\n\n[Composer\\Downloader\\TransportException]\ncurl error 6 while downloading https://flex.symfony.com/versions.json: Could not resolve host: flex.symfony.com",
			doc:   "We weren't able to retrieve one or more dependencies. Please check your Internet connection and try again.",
		},
		{
			error: "The following exception probably indicates you are offline or have misconfigured DNS resolver(s)\n\n[Composer\\Downloader\\TransportException]\ncurl error 6 while downloading https://flex.symfony.com/versions.json: Could not resolve host: flex.symfony.com",
			doc:   "We weren't able to retrieve one or more dependencies. Please check your Internet connection and try again.",
		},
		{
			error: "Root composer.json requires drupal/entity_pager 1.0@RC, found drupal/entity_pager[dev-1.x, dev-2.0.x, 1.0.0-alpha1, ..., 1.x-dev (alias of dev-1.x), 2.0.x-dev (alias of dev-2.0.x)] but it does not match the constraint.",
			doc:   "Couldn't resolve version drupal/entity_pager 1.0@RC , please make sure it exists:\n",
		},
		{
			error: "Loading composer repositories with package information\nUpdating dependencies\nYour requirements could not be resolved to an installable set of packages.\n\n  Problem 1\n    - Root composer.json requires blablabla/blabla, it could not be found in any version, there may be a typo in the package name.\n\nPotential causes:\n - A typo in the package name\n - The package is not available in a stable-enough version according to your minimum-stability setting\n   see <https://getcomposer.org/doc/04-schema.md#minimum-stability> for more details.\n - It's a private package and you forgot to add a custom repository to find it\n\nRead <https://getcomposer.org/doc/articles/troubleshooting.md> for further common problems.\n",
			doc:   "An error occurred during dependencies resolve for: blablabla/blabla\n\nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.\n\n",
		},
	}

	for _, c := range cases {
		expectedError := util.NewPMJobError(c.error)
		expectedError.SetDocumentation(c.doc)

		cmdErr := errors.New(c.error)
		j := NewJob("file", true, testdata.CmdFactoryMock{InstallCmdName: "echo", MakeInstallErr: cmdErr})

		go jobTestdata.WaitStatus(j)

		j.Run()

		errors := j.Errors().GetAll()

		assert.Len(t, errors, 1)
		assert.Contains(t, errors, expectedError)
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
