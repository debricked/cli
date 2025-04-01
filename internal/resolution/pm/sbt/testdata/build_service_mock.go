package testdata

type BuildServiceMock struct {
	Value []string
	Err   error
}

func (b BuildServiceMock) ParseBuildModules(_ string) ([]string, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	return b.Value, nil
}
