package sbt

import (
	"os"
	"path/filepath"
	"regexp"
)

type IBuildService interface {
	ParseBuildModules(path string) ([]string, error)
	FindPomFile(dir string) (string, error)
	RenamePomToXml(pomFile, destDir string) (string, error)
}

type BuildService struct{}

func (b BuildService) ParseBuildModules(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	moduleRegex := regexp.MustCompile(`project\s*\(\s*"([^"]+)"\s*\)`)
	matches := moduleRegex.FindAllStringSubmatch(string(content), -1)

	modules := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			modules = append(modules, match[1])
		}
	}

	return modules, nil
}

func (b BuildService) FindPomFile(dir string) (string, error) {
	targetDir := filepath.Join(dir, "target")

	scalaVersionDirs, err := filepath.Glob(filepath.Join(targetDir, "scala-*"))
	if err != nil || len(scalaVersionDirs) == 0 {
		return "", err
	}

	for _, scalaDir := range scalaVersionDirs {
		pomFiles, err := filepath.Glob(filepath.Join(scalaDir, "*.pom"))
		if err == nil && len(pomFiles) > 0 {
			return pomFiles[0], nil
		}
	}

	return "", nil
}

func (b BuildService) RenamePomToXml(pomFile, destDir string) (string, error) {
	content, err := os.ReadFile(pomFile)
	if err != nil {
		return "", err
	}

	pomXmlPath := filepath.Join(destDir, "pom.xml")
	err = os.WriteFile(pomXmlPath, content, 0600)
	if err != nil {
		return "", err
	}

	return pomXmlPath, nil
}
