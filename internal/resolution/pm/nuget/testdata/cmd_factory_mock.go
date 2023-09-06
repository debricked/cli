package testdata

import (
	"os/exec"
)

type CmdFactoryMock struct {
	InstallCmdName string
	MakeInstallErr error
	CmdOutput      []byte
	CmdError       error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		InstallCmdName: "echo",
	}
}
func (f CmdFactoryMock) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	if f.MakeInstallErr != nil {
		return nil, f.MakeInstallErr
	}

	cmd := exec.Command(f.InstallCmdName)
	if f.CmdError != nil {
		cmd = exec.Command("sh", "-c", f.InstallCmdName+" && exit 1")
	}

	if f.CmdOutput != nil {
		cmd = exec.Command("sh", "-c", "echo -n \""+string(f.CmdOutput)+"\" && exit 1")
	}

	return cmd, nil
}
