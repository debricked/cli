package golanfinder

import (
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/debricked/cli/internal/file"
)

type GolangFinder struct{}

func (f GolangFinder) FindRoots(files []string) ([]string, error) {
	mainFiles := finder.FilterFiles(files, "main.go")
	return mainFiles, nil
}

func (f GolangFinder) FindDependencyDirs(files []string, findJars bool) ([]string, error) {
	// Not needed for golang
	return []string{}, nil
}

func (f GolangFinder) FindFiles(roots []string, exclusions []string) ([]string, error) {
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

			if !info.IsDir() && !excluded && filepath.Ext(path) == ".go" {
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
