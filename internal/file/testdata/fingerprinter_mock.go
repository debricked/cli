package testdata

import (
	"github.com/debricked/cli/internal/file"
)

type FingerprintMock struct {
	error error
}

func NewFingerprintMock() *FingerprintMock {
	return &FingerprintMock{
		error: nil,
	}
}

func (f *FingerprintMock) FingerprintFiles(rootPath string, exclusions []string) (file.Fingerprints, error) {
	return file.Fingerprints{}, f.error
}
