package finder

type IFinder interface {
	FindRoots(files []string) ([]string, error)
	FindDependencyDirs(files []string, findJars bool) ([]string, error)
	FindFiles(paths []string, exclusions []string) ([]string, error)
}
