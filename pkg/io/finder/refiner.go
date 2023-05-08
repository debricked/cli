package finder

import (
	"os"
	"path/filepath"
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

func MapFilesToDir(dirs []string, files []string) map[string][]string {
	dirToFilesMap := make(map[string][]string)

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
