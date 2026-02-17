package testdata

import "os/exec"

type CmdFactoryMock struct {
	InstallCmdName string
	MakeInstallErr error
}

type ExecPathMock struct{}

func (ExecPathMock) LookPath(file string) (string, error) {
	// Just return the name; Exec.Cmd will use PATH resolution
	return file, nil
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		InstallCmdName: "echo",
	}
}

func (f CmdFactoryMock) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	return exec.Command(f.InstallCmdName), f.MakeInstallErr
}
