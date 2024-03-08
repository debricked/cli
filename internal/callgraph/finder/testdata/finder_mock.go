package testdata

type FinderMock struct {
	FindDependencyDirsNames []string
	FindDependencyDirsErr   error
	FindRootsNames          []string
	FindRootsErr            error
	FindFilesNames          []string
	FindFilesErr            error
}

func NewEmptyFinderMock() FinderMock {
	return FinderMock{
		FindDependencyDirsNames: []string{},
		FindRootsNames:          []string{},
		FindFilesNames:          []string{},
	}
}

func (f FinderMock) FindDependencyDirs(_ []string, _ bool) ([]string, error) {
	return f.FindDependencyDirsNames, f.FindDependencyDirsErr
}

func (f FinderMock) FindRoots(_ []string) ([]string, error) {
	return f.FindRootsNames, f.FindRootsErr
}

func (f FinderMock) FindFiles(_ []string, _ []string) ([]string, error) {
	return f.FindFilesNames, f.FindFilesErr
}
