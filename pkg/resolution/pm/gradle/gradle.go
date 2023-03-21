package gradle

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type IGradleService interface {
	GetRootGradlePath(files []string) []string
}

type GradleService struct {
}

func ParseJavaGradleFile(file string) []string {
	includeRegex := regexp.MustCompile(`include\(\".*\"\)`)
	matches := includeRegex.FindAllString(file, -1)
	subProjects := make([]string, 0)
	for _, match := range matches {
		match = strings.Replace(match, "include(", "", -1)
		match = strings.Replace(match, ")", "", -1)
		match = strings.Replace(match, "\"", "", -1)
		match = strings.Replace(match, " ", "", -1)
		subProjects = append(subProjects, match)
	}
	return subProjects
}

func ParseKotlinGradleFile(filepath string) ([]string, error) {
	includeRegex := regexp.MustCompile(`include\(\"([^\"]+)\"`)
	subProjects := make([]string, 0)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	matches := includeRegex.FindAllString(string(content), -1)
	for _, match := range matches {
		match = strings.Replace(match, "include(", "", -1)
		match = strings.Replace(match, ")", "", -1)
		match = strings.Replace(match, "\"", "", -1)
		match = strings.Replace(match, " ", "", -1)
		subProjects = append(subProjects, match)
	}
	return subProjects, nil
}

func (t GradleService) GetRootGradlePath(files []string) []string {
	settingsMap := make(map[string]bool)
	buildMap := make(map[string]bool)
	childMap := make(map[string]bool)
	rootGradlePaths := make([]string, 0)
	validGradlePaths := make([]string, 0)
	for _, file := range files {
		if filepath.Base(file) == "settings.gradle" || filepath.Base(file) == "settings.gradle.kts" {
			settingsMap[file] = true
		} else if filepath.Base(file) == "build.gradle" || filepath.Base(file) == "build.gradle.kts" {
			buildMap[file] = true
			validGradlePaths = append(validGradlePaths, file)
		}
	}

	for file := range settingsMap {
		var childs []string
		var err error
		var rootName string
		extension := filepath.Ext(file)
		if extension == ".kts" {
			childs, err = ParseKotlinGradleFile(file)
			if err != nil {
				continue
			}
			rootName = "build.gradle.kts"
		} else {
			childs = ParseJavaGradleFile(file)
			rootName = "build.gradle"
		}

		if len(childs) == 0 {
			continue
		}

		for _, child := range childs {
			childPath := filepath.Join(filepath.Dir(file), "subprojects", filepath.Dir(child), filepath.Base(child), rootName)
			if _, ok := buildMap[childPath]; ok {
				childMap[childPath] = true // add child if setting has a sibling build ifle
			}
		}
	}

	for _, file := range validGradlePaths {
		if _, ok := childMap[file]; !ok {
			rootGradlePaths = append(rootGradlePaths, file)
		}
	}

	return rootGradlePaths
}
