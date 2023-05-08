package finder

import "path/filepath"

func FindJavaClassDirs(files []string) []string {
	filteredFiles := FilterFiles(files, "*.class")
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
