package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FilterFiles(files []string, pattern string) []string {
	filteredFiles := []string{}
	for _, file := range files {
		matched, _ := filepath.Match(pattern, filepath.Base(file))
		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func ConvertPathsToAbsPaths(paths []string) ([]string, error) {
	absPaths := []string{}

	for _, path := range paths {
		path, err := filepath.Abs(path)

		if err != nil {
			return []string{}, err
		}

		absPaths = append(absPaths, path)
	}

	return absPaths, nil
}

// Matches class directories to closest root pom file and creates a map
// with each root pom file pointing at a list of its related class directories
func MapFilesToDir(rootPomFiles []string, classDirs []string) map[string][]string {
	pomFileToClassDirsMap := make(map[string][]string)

	if len(rootPomFiles) == 0 {
		return pomFileToClassDirsMap
	}

	for _, classDir := range classDirs {
		matchingPomFile, err := findPomFileMatch(classDir, rootPomFiles)
		if err != nil {
			continue
		}

		if _, ok := pomFileToClassDirsMap[matchingPomFile]; !ok {
			pomFileToClassDirsMap[matchingPomFile] = []string{}
		}
		pomFileToClassDirsMap[matchingPomFile] = append(pomFileToClassDirsMap[matchingPomFile], classDir)
	}

	return pomFileToClassDirsMap
}

func findPomFileMatch(classDir string, pomFiles []string) (string, error) {
	var matchingPomFile string
	longestSeperatorMatch := 0
	matched := false

	for _, pomFile := range pomFiles {
		pomFilePath := strings.TrimSuffix(pomFile, "pom.xml")
		if strings.Contains(classDir, pomFilePath) {
			numberSeparators := strings.Count(pomFilePath, string(os.PathSeparator))
			if numberSeparators > longestSeperatorMatch {
				matched = true
				longestSeperatorMatch = numberSeparators
				matchingPomFile = pomFile
			}
		}
	}

	if !matched {
		return "", fmt.Errorf("No part of the path matches")
	}

	return matchingPomFile, nil
}

func GCDPath(paths []string) string {
	var result string
	var shortest string

	for i, path := range paths {
		if i == 0 || len(path) < len(shortest) {
			shortest = path
		}
	}

	for i := 0; i < len(shortest); i++ {
		c := shortest[i]

		if filepath.Separator == c {
			dirpath := shortest[:i+1]
			for _, path := range paths {
				if !strings.HasPrefix(path, dirpath) {
					return result
				}
			}

			result = dirpath
		}
	}

	return result
}
