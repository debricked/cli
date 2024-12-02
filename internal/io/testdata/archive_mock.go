package testdata

import "strings"

type ArchiveMock struct {
	ZipFileError   error
	B64Error       error
	CleanupError   error
	PathError      error
	Dir            string
	UnzipFileError error
}

func (am ArchiveMock) ZipFile(sourceName string, targetName string, zipName string) error {
	if !strings.HasPrefix(sourceName, am.Dir) || !strings.HasPrefix(targetName, am.Dir) {
		return am.PathError
	}
	return am.ZipFileError

}

func (am ArchiveMock) B64(sourceName string, targetName string) error {
	return am.B64Error

}

func (am ArchiveMock) Cleanup(fileName string) error {
	return am.CleanupError
}

func (am ArchiveMock) UnzipFile(sourcePath string, targetPath string) error {
	return am.UnzipFileError
}
