package java

import (
	"errors"
	"fmt"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	ioWriter "github.com/debricked/cli/pkg/io/writer"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	config     conf.IConfig
}

func NewJob(files []string, cmdFactory ICmdFactory, writer ioWriter.IFileWriter, config conf.IConfig) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(files),
		cmdFactory: cmdFactory,
		config:     config,
	}
}

func (j *Job) Run() {
	fmt.Println("ENTERED RUN")
	workingDirectory := "." // filepath.Dir(filepath.Clean(j.GetFile()))
	targetClasses := "/home/magnus/Projects/exploration/dependency-demo-app/target/classes/"
	dependencyClasses := "/home/magnus/Projects/exploration/dependency-demo-app/target/dependency/"
	cmd, err := j.cmdFactory.MakeCallGraphGenerationCmd(workingDirectory, targetClasses, dependencyClasses)
	if err != nil {
		j.Errors().Critical(err)

		return
	}
	j.SendStatus("creating dependency graph")
	var output []byte
	fmt.Println("run command", cmd.Args)
	output, err = cmd.Output()
	fmt.Println("done running command", cmd.Args)
	if err != nil {
		if output == nil {
			j.Errors().Critical(err)
		} else {
			j.Errors().Critical(errors.New(string(output)))
		}
	}
}
