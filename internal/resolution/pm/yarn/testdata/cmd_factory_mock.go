package testdata

import (
	"os/exec"
)

type CmdFactoryMock struct {
	InstallCmdName string
	MakeInstallErr error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		// Use the Go binary, which is guaranteed to be available
		// wherever `go test` is running (both locally and in CI).
		InstallCmdName: "go",
	}
}

func (f CmdFactoryMock) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	// When using the default mock command, run a harmless
	// `go version` invocation to avoid depending on platform-
	// specific binaries like `echo`.
	if f.InstallCmdName == "" || f.InstallCmdName == "go" {
		return exec.Command("go", "version"), f.MakeInstallErr
	}

	// For tests that explicitly set InstallCmdName (e.g. to a
	// bad value), keep that behavior so error paths are exercised.
	return exec.Command(f.InstallCmdName), f.MakeInstallErr
}
