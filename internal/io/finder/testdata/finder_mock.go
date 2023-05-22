package testdata

type FinderMock struct {
	FindJavaClassDirsNames []string
	FindJavaClassDirsErr   error
	FindMavenRootsNames    []string
	FindMavenRootsErr      error
	FindFilesNames         []string
	FindFilesErr           error
}

func NewEmptyFinderMock() FinderMock {
	return FinderMock{
		FindJavaClassDirsNames: []string{},
		FindMavenRootsNames:    []string{},
		FindFilesNames:         []string{},
	}
}

func (f FinderMock) FindJavaClassDirs(_ []string) ([]string, error) {
	return f.FindJavaClassDirsNames, f.FindJavaClassDirsErr
}

func (f FinderMock) FindMavenRoots(_ []string) ([]string, error) {
	return f.FindMavenRootsNames, f.FindMavenRootsErr
}

func (f FinderMock) FindFiles(_ []string, _ []string) ([]string, error) {
	return f.FindFilesNames, f.FindFilesErr
}
