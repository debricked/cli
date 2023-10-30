package finder

import (
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/finder/maven"
	"github.com/debricked/cli/internal/file"
)

type IFinder interface {
	FindMavenRoots(files []string) ([]string, error)
	FindJavaClassDirs(files []string, findJars bool) ([]string, error)
	FindFiles(paths []string, exclusions []string) ([]string, error)
}

type Finder struct{}

func (f Finder) FindMavenRoots(files []string) ([]string, error) {
	pomFiles := FilterFiles(files, "pom.xml")
	ps := maven.PomService{}
	rootFiles := ps.GetRootPomFiles(pomFiles)

	return rootFiles, nil
}

func (f Finder) FindJavaClassDirs(files []string, findJars bool) ([]string, error) {
	filteredFiles := FilterFiles(files, ".*\\.class")
	dirsWithClassFiles := make(map[string]bool)
	for _, file := range filteredFiles {
		dirsWithClassFiles[filepath.Dir(file)] = true
	}

	dirJarFiles := []string{}
	for key := range dirsWithClassFiles {
		dirJarFiles = append(dirJarFiles, key)
	}

	if findJars {
		filteredJarFiles := FilterFiles(files, ".*\\.jar")
		dirJarFiles = append(dirJarFiles, filteredJarFiles...)
	}

	return dirJarFiles, nil
}

func (f Finder) FindFiles(roots []string, exclusions []string) ([]string, error) {
	files := make(map[string]bool)
	var err error = nil

	for _, root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			excluded := file.Excluded(exclusions, path)

			if info.IsDir() && excluded {
				return filepath.SkipDir
			}

			if !info.IsDir() && !excluded {
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
