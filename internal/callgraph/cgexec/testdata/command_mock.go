package testdata

import "github.com/stretchr/testify/mock"

type CmdMock struct {
	mock.Mock
}

func NewCmdMock() CmdMock {
	return CmdMock{}
}

func (m *CmdMock) CombinedOutput() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *CmdMock) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *CmdMock) Wait() error {
	args := m.Called()
	return args.Error(0)
}
