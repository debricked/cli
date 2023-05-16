package testdata

import (
	"github.com/debricked/cli/internal/resolution/file"
	"github.com/debricked/cli/internal/resolution/pm"
)

type BatchFactoryMock struct {
	pms []pm.IPm
}

func NewBatchFactoryMock() BatchFactoryMock {
	return BatchFactoryMock{
		pms: pm.Pms(),
	}
}

func (bf BatchFactoryMock) Make(_ []string) []file.IBatch {

	return []file.IBatch{}
}
