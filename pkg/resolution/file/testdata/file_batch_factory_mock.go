package testdata

import (
	"github.com/debricked/cli/pkg/resolution/file"
	"github.com/debricked/cli/pkg/resolution/pm"
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
