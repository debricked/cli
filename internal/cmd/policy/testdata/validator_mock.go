package testdata

import (
	"github.com/debricked/cli/internal/policy"
)

type ValidatorMock struct {
	Path   string
	Result policy.Result
	Error  error
}

func (mock *ValidatorMock) Validate(path string) (policy.Result, error) {
	mock.Path = path

	return mock.Result, mock.Error
}
