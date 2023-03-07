package maven

import (
	"path/filepath"

	"github.com/vifraa/gopom"
)

type IPomX interface {
	GetRootPomFiles(files []string) []string
	ParsePomModules(path string) ([]string, error)
}

type PomX struct{}

func (p PomX) ParsePomModules(path string) ([]string, error) {

	pom, err := gopom.Parse(path)

	if err != nil {
		return nil, err
	}

	return pom.Modules, nil
}

func (p PomX) GetRootPomFiles(files []string) []string {

	childMap := make(map[string][]string)

	for _, file_path := range files {

		modules, _ := p.ParsePomModules(file_path)

		if len(modules) == 0 {
			continue
		}

		for _, module := range modules {

			// path to child pom
			modulePath := filepath.Join(filepath.Dir(file_path), filepath.Dir(module), filepath.Base(module), "pom.xml")

			childMap[modulePath] = append(childMap[modulePath], file_path)

		}
	}

	roots := make([]string, 0)

	for _, file := range files {
		if _, ok := childMap[file]; !ok {
			roots = append(roots, file)
		}
	}

	return roots

}
