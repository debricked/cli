package testdata

import (
	"bytes"
	"errors"
	"os"
)

type CmdConfig struct {
	Start error
}

func NewCmdConfig() *CmdConfig {
	return &CmdConfig{errors.New("test error")}
}

type CommandMock struct {
	cmdConfig CmdConfig
}

func NewCommandMock() CommandMock {
	return CommandMock{}
}

func NewCommandMockWithConfig(cmdConfig CmdConfig) CommandMock {
	return CommandMock{cmdConfig}
}

func (m CommandMock) CombinedOutput() ([]byte, error) {
	return []byte{}, nil
}

func (m CommandMock) Start() error {
	return m.cmdConfig.Start
}

func (m CommandMock) Wait() error {
	return nil
}

func (m CommandMock) GetProcess() *os.Process {
	return nil
}

func (m CommandMock) SetStderr(stderr *bytes.Buffer) {
}

func (m CommandMock) SetStdout(stdout *bytes.Buffer) {
}

func (m CommandMock) GetArgs() []string {
	return []string{"mvn", "package", "TEST"}
}

func (m CommandMock) GetDir() string {
	return "."
}
