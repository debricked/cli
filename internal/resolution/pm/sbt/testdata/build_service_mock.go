package testdata

type BuildServiceMock struct {
	Value []string
	Err   error
}

func (b BuildServiceMock) ParseBuildModules(_ string) ([]string, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	if b.Value == nil {
		return []string{"default-module"}, nil
	}

	return b.Value, nil
}

func (b BuildServiceMock) FindPomFile(_ string) (string, error) {
	if b.Err != nil {
		return "", b.Err
	}

	return "pom.xml", nil
}

func (b BuildServiceMock) RenamePomToXml(pomFile, destDir string) (string, error) {
	if b.Err != nil {
		return "", b.Err
	}

	return pomFile, nil
}
