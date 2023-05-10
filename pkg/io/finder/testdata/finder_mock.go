package testdata

type FinderMock struct {
	FindJavaClassDirsNames []string
	FindJavaClassDirsErr   error
	FindMavenRootsNames    []string
	FindMavenRootsErr      error
	FindGradleRootsNames   []string
	FindGradleRootsErr     error
}

func NewEmptyFinderMock() FinderMock {
	return FinderMock{
		FindJavaClassDirsNames: []string{},
		FindMavenRootsNames:    []string{},
		FindGradleRootsNames:   []string{},
	}
}

func (f FinderMock) FindJavaClassDirs(_ []string) ([]string, error) {
	return f.FindJavaClassDirsNames, f.FindJavaClassDirsErr
}

func (f FinderMock) FindMavenRoots(_ []string) ([]string, error) {
	return f.FindMavenRootsNames, f.FindMavenRootsErr
}

func (f FinderMock) FindGradleRoots(_ []string) ([]string, error) {
	return f.FindGradleRootsNames, f.FindGradleRootsErr
}
