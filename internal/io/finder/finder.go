package finder

import (
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/io/finder/maven"
)

type IFinder interface {
	FindMavenRoots(files []string) ([]string, error)
	FindJavaClassDirs(files []string) ([]string, error)
	FindFiles(paths []string, exclusions []string) ([]string, error)
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

func (f Finder) FindFiles(roots []string, exclusions []string) ([]string, error) {
	files := make(map[string]bool)
	var err error = nil

	for _, root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, dir := range exclusions {
				if info.IsDir() && info.Name() == dir {
					return filepath.SkipDir
				}
			}

			if !info.IsDir() {
				files[path] = true
			}

			return nil
		})

		if err != nil {
			break
		}
	}

	fileList := make([]string, len(files))
	i := 0
	for k := range files {
		fileList[i] = k
		i++
	}

	return fileList, err
}
