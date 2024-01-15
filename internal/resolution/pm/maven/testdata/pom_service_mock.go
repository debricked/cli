package testdata

type PomServiceMock struct {
	Value []string
	Err   error
}

func (p PomServiceMock) ParsePomModules(_ string) ([]string, error) {
	if p.Err != nil {
		return nil, p.Err
	}

	return p.Value, nil
}
