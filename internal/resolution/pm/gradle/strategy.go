package gradle

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

type Strategy struct {
	files       []string
	paths       []string
	ErrorWriter io.Writer
	GradleSetup ISetup
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	fileWriter := writer.FileWriter{}
	factory := CmdFactory{}
	gradleSetup, err := s.GradleSetup.Configure(s.files, s.paths)
	if err != nil {
		if _, ok := err.(SetupSubprojectError); ok {
			warningColor := color.New(color.FgYellow, color.Bold).SprintFunc()
			defaultOutputWriter := log.Writer()
			log.SetOutput(s.ErrorWriter)
			log.Println(warningColor("Warning:\n") + err.Error())
			log.SetOutput(defaultOutputWriter)
		} else {
			return nil, err
		}
	}
	gradleMainDirs := make(map[string]bool)
	for _, gradleProject := range gradleSetup.GradleProjects {
		dir := gradleProject.dir
		if _, ok := gradleMainDirs[dir]; ok {
			continue
		}
		gradleMainDirs[dir] = true
		jobs = append(jobs, NewJob(gradleProject.mainBuildFile, dir, gradleProject.gradlew, gradleSetup.groovyScriptPath, factory, fileWriter))

	}
	for _, file := range s.files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		if _, ok := gradleSetup.subProjectMap[dir]; ok {
			continue
		}
		if _, ok := gradleMainDirs[dir]; ok {
			continue
		}
		gradleMainDirs[dir] = true
		gradlew := gradleSetup.GetGradleW(dir)
		jobs = append(jobs, NewJob(file, dir, gradlew, gradleSetup.groovyScriptPath, factory, fileWriter))
	}

	return jobs, nil
}

func NewStrategy(files []string, paths []string) Strategy {
	return Strategy{files, paths, os.Stdout, NewGradleSetup()}
}
