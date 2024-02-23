package finder

type IFinder interface {
	FindMavenRoots(files []string) ([]string, error)
	FindJavaClassDirs(files []string, findJars bool) ([]string, error)
	FindFiles(paths []string, exclusions []string) ([]string, error)
}
