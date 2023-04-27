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
	var validFiles []string
	var roots []string

	for _, filePath := range files {
		modules, err := p.ParsePomModules(filePath)

		if err != nil {
			continue
		}

		validFiles = append(validFiles, filePath)

		if len(modules) == 0 {
			continue
		}

		for _, module := range modules {
			modulePath := filepath.Join(filepath.Dir(filePath), filepath.Dir(module), filepath.Base(module), "pom.xml")
			childMap[modulePath] = true
		}
	}

	for _, file := range validFiles {
		if _, ok := childMap[file]; !ok {
			roots = append(roots, file)
		}
	}

	return roots
}
