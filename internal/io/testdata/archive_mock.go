package testdata

type ArchiveMock struct {
	ZipFileError error
	B64Error     error
	CleanupError error
}

func (am ArchiveMock) ZipFile(sourceName string, targetName string) error {
	return am.ZipFileError

}

func (am ArchiveMock) B64(sourceName string, targetName string) error {
	return am.B64Error

}

func (am ArchiveMock) Cleanup(fileName string) error {
	return am.CleanupError
}
