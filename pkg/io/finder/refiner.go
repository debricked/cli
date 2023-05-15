package finder

import (
	"os"
	"path/filepath"
	"strings"
)

func FindFiles(roots []string, exclusions []string) ([]string, error) {
	files := make(map[string]bool)
	var err error = nil

	for _, root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, dir := range exclusions {
				if info.IsDir() && info.Name() == dir {
					return filepath.SkipDir
				}
			}

			if !info.IsDir() {
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

func MapFilesToDir(dirs []string, files []string) map[string][]string {
	dirToFilesMap := make(map[string][]string)

	if len(dirs) == 0 {
		return dirToFilesMap
	}

	for _, file := range files {
		longestMatchLength := 0
		var matchingDir string
		for _, dir := range dirs {
			matchLength := 0
			for i := 0; i < len(file) && i < len(dir); i++ {
				if file[i] != dir[i] {
					break
				}
				matchLength++
			}
			if matchLength > longestMatchLength {
				longestMatchLength = matchLength
				matchingDir = dir
			}
		}

		if _, ok := dirToFilesMap[matchingDir]; ok == false {
			dirToFilesMap[matchingDir] = []string{}
		}
		dirToFilesMap[matchingDir] = append(dirToFilesMap[matchingDir], file)
	}

	return dirToFilesMap
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
