package gradle

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type Strategy struct {
	files       []string
	paths       []string
	ErrorWriter io.Writer
	GradleSetup IGradleSetup
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	writer := writer.FileWriter{}
	factory := CmdFactory{}
	gradleSetup, err := s.GradleSetup.Setup(s.files, s.paths)

	if err != nil {

		if _, ok := err.(GradleSetupSubprojectError); ok {
			warningColor := color.New(color.FgYellow, color.Bold).SprintFunc()
			defaultOutputWriter := log.Writer()
			log.SetOutput(s.ErrorWriter)
			log.Println(warningColor("Warning:\n") + err.Error())
			log.SetOutput(defaultOutputWriter)
		} else {
			return nil, err
		}
	}

	for _, gradleProject := range gradleSetup.GradleProjects {
		jobs = append(jobs, NewJob(gradleProject.dir, gradleProject.gradlew, gradleSetup.groovyScriptPath, factory, writer))
		fmt.Println("Added job for " + gradleProject.dir)
	}
	fmt.Println(s.files)
	for _, file := range s.files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		fmt.Println("Found dir" + dir)
		if _, ok := gradleSetup.subProjectMap[dir]; ok {
			continue
		}
		gradlew := gradleSetup.GetGradleW(dir)
		jobs = append(jobs, NewJob(dir, gradlew, gradleSetup.groovyScriptPath, factory, writer))
	}
	return jobs, nil
}

func NewStrategy(files []string, paths []string) Strategy {
	return Strategy{files, paths, os.Stdout, NewGradleSetup()}
}
