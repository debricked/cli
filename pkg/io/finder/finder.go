package finder

import (
	"path/filepath"

	"github.com/debricked/cli/pkg/io/finder/gradle"
	"github.com/debricked/cli/pkg/io/finder/maven"
)

type IFinder interface {
	FindMavenRoots(files []string) ([]string, error)
	FindJavaClassDirs(files []string) ([]string, error)
	FindGradleRoots(files []string) ([]string, error)
}

type Finder struct{}

func (f Finder) FindMavenRoots(files []string) ([]string, error) {
	pomFiles := FilterFiles(files, "pom.xml")
	ps := maven.PomService{}
	rootFiles := ps.GetRootPomFiles(pomFiles)
	return rootFiles, nil
}

func (f Finder) FindJavaClassDirs(files []string) ([]string, error) {
	filteredFiles := FilterFiles(files, "*.class")
	dirsWithJarFiles := make(map[string]bool)
	for _, file := range filteredFiles {
		dirsWithJarFiles[filepath.Dir(file)] = true
	}

	jarFiles := []string{}
	for key := range dirsWithJarFiles {
		jarFiles = append(jarFiles, key)
	}

	return jarFiles, nil
}

func (f Finder) FindGradleRoots(files []string) ([]string, error) {
	gradleBuildFiles := FilterFiles(files, "gradle.build(.kts)?")
	gradleSetup := gradle.NewGradleSetup()
	err := gradleSetup.Configure(files)
	if err != nil {

		return []string{}, err
	}

	gradleMainDirs := make(map[string]bool)
	for _, gradleProject := range gradleSetup.GradleProjects {
		dir := gradleProject.Dir
		if _, ok := gradleMainDirs[dir]; ok {
			continue
		}
		gradleMainDirs[dir] = true
	}
	for _, file := range gradleBuildFiles {
		dir, _ := filepath.Abs(filepath.Dir(file))
		if _, ok := gradleSetup.SubProjectMap[dir]; ok {
			continue
		}
		if _, ok := gradleMainDirs[dir]; ok {
			continue
		}
		gradleMainDirs[dir] = true
	}

	roots := []string{}
	for key := range gradleMainDirs {
		roots = append(roots, key)
	}

	return roots, nil
}
