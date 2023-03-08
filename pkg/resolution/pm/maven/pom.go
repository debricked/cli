package maven

import (
	"path/filepath"

	"github.com/vifraa/gopom"
)

type IPomService interface {
	GetRootPomFiles(files []string) []string
	ParsePomModules(path string) ([]string, error)
}

type PomService struct{}

func (p PomService) ParsePomModules(path string) ([]string, error) {
	pom, err := gopom.Parse(path)

	if err != nil {
		return nil, err
	}

	return pom.Modules, nil
}

func (p PomService) GetRootPomFiles(files []string) []string {
	childMap := make(map[string]bool)
	validFiles := make([]string, 0)
	roots := make([]string, 0)

	for _, file_path := range files {
		modules, err := p.ParsePomModules(file_path)

		if err != nil {
			continue
		}

		validFiles = append(validFiles, file_path)

		if len(modules) == 0 {
			continue
		}

		for _, module := range modules {
			modulePath := filepath.Join(filepath.Dir(file_path), filepath.Dir(module), filepath.Base(module), "pom.xml")
			childMap[modulePath] = bool(true)
		}
	}

	for _, file := range validFiles {
		if _, ok := childMap[file]; !ok {
			roots = append(roots, file)
		}
	}

	return roots
}
