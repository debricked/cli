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
	cmdConfig           CmdConfig
	Process             *os.Process
	CombinedOutputError error
	SignalError         error
	WaitError           error
}

func NewCommandMock() CommandMock {
	return CommandMock{}
}

func NewCommandMockWithConfig(cmdConfig CmdConfig) CommandMock {
	return CommandMock{cmdConfig, nil, nil, nil, nil}
}

func (m CommandMock) CombinedOutput() ([]byte, error) {
	return []byte{}, m.CombinedOutputError
}

func (m CommandMock) Start() error {
	return m.cmdConfig.Start
}

func (m CommandMock) Wait() error {
	return m.WaitError
}

func (m CommandMock) GetProcess() *os.Process {
	return m.Process
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

func (m CommandMock) Signal(process *os.Process, signal os.Signal) error {
	return m.SignalError
}

func (m CommandMock) GetStdOut() *bytes.Buffer {
	return &bytes.Buffer{}
}

func (m CommandMock) GetStdErr() *bytes.Buffer {
	return &bytes.Buffer{}
}
