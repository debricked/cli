package gradle

import (
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type Strategy struct {
	files []string
	paths []string
}

func (s Strategy) Invoke() []job.IJob {
	var jobs []job.IJob

	writer := writer.FileWriter{}
	factory := CmdFactory{}
	gradleSetup := NewGradleSetup()
	gradleSetup.Setup(s.files, s.paths)

	for _, gradleProject := range gradleSetup.GradleProjects {
		jobs = append(jobs, NewJob(gradleProject.dir, gradleProject.gradlew, gradleSetup.groovyScriptPath, factory, writer))
	}

	for _, file := range s.files {
		dir, _ := filepath.Abs(filepath.Dir(file))
		if _, ok := gradleSetup.subProjectMap[dir]; ok {
			continue
		}

		gradlew := gradleSetup.GetGradleW(dir)
		jobs = append(jobs, NewJob(dir, gradlew, gradleSetup.groovyScriptPath, factory, writer))
	}

	// Add teardown (remove initfiles ?)
	fmt.Println("SubProjects found", len(gradleSetup.subProjectMap))
	fmt.Println("gradlew found", len(gradleSetup.gradlewMap))
	fmt.Println("Jobs found", len(jobs))

	return jobs
}

func NewStrategy(files []string, paths []string) Strategy {
	return Strategy{files, paths}
}
