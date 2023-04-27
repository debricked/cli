package finder

import (
	"path/filepath"

	"github.com/debricked/cli/internal/io/finder/maven"
)

type IFinder interface {
	FindMavenRoots(files []string) ([]string, error)
	FindJavaClassDirs(files []string) ([]string, error)
}

type Finder struct{}

func (f Finder) FindMavenRoots(files []string) ([]string, error) {
	pomFiles := FilterFiles(files, "pom.xml")
	ps := maven.PomService{}
	rootFiles := ps.GetRootPomFiles(pomFiles)

	return rootFiles, nil
}

func (f Finder) FindJavaClassDirs(files []string) ([]string, error) {
	filteredFiles := FilterFiles(files, "*.class")
	dirsWithJarFiles := make(map[string]bool)
	for _, file := range filteredFiles {
		dirsWithJarFiles[filepath.Dir(file)] = true
	}

	jarFiles := []string{}
	for key := range dirsWithJarFiles {
		jarFiles = append(jarFiles, key)
	}

	return jarFiles, nil
}
