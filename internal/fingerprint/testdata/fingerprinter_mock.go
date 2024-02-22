package testdata

import (
	"github.com/debricked/cli/internal/fingerprint"
)

type FingerprintMock struct {
	error error
}

func NewFingerprintMock() *FingerprintMock {
	return &FingerprintMock{
		error: nil,
	}
}

func (f *FingerprintMock) FingerprintFiles(rootPath string, exclusions []string, fingerprintCompressedContent bool, includeSnippet bool) (fingerprint.Fingerprints, error) {
	return fingerprint.Fingerprints{}, f.error
}
