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

func NewFingerprintMockFileExistsError() *FingerprintMock {
	return &FingerprintMock{
		error: &fingerprint.FingerprintFileExistsError{},
	}
}

func (f *FingerprintMock) FingerprintFiles(
	options fingerprint.DebrickedOptions,
) (fingerprint.Fingerprints, error) {
	return fingerprint.Fingerprints{}, f.error
}
