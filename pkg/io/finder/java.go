package finder

import "path/filepath"

func FindJarDirs(files []string) []string {
	filteredFiles := FilterFiles(files, "*.jar")
	dirsWithJarFiles := make(map[string]bool)
	for _, file := range filteredFiles {
		dirsWithJarFiles[filepath.Dir(file)] = true
	}

	jarFiles := []string{}
	for key := range dirsWithJarFiles {
		jarFiles = append(jarFiles, key)
	}

	return jarFiles
}
