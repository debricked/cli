package javafinder

import (
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/debricked/cli/internal/file"
)

type JavaFinder struct{}

func (f JavaFinder) FindRoots(files []string) ([]string, error) {
	pomFiles := finder.FilterFiles(files, "pom.xml")
	ps := PomService{}
	rootFiles := ps.GetRootPomFiles(pomFiles)

	return rootFiles, nil
}

func (f JavaFinder) FindDependencyDirs(files []string, findJars bool) ([]string, error) {
	filteredFiles := finder.FilterFiles(files, ".*\\.class")
	dirsWithClassFiles := make(map[string]bool)
	for _, file := range filteredFiles {
		dirsWithClassFiles[filepath.Dir(file)] = true
	}

	dirJarFiles := []string{}
	for key := range dirsWithClassFiles {
		dirJarFiles = append(dirJarFiles, key)
	}

	if findJars {
		filteredJarFiles := finder.FilterFiles(files, ".*\\.jar")
		dirJarFiles = append(dirJarFiles, filteredJarFiles...)
	}

	return dirJarFiles, nil
}

func (f JavaFinder) FindFiles(roots []string, exclusions []string) ([]string, error) {
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
