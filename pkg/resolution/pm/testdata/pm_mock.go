package testdata

type PmMock struct {
	N  string
	Ms []string
}

func (pm PmMock) Name() string {
	return pm.N
}

func (pm PmMock) Manifests() []string {
	return pm.Ms
}
