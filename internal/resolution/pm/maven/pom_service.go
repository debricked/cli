package maven

import (
	"github.com/vifraa/gopom"
)

type IPomService interface {
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
