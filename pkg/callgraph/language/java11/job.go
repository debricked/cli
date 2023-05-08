package java

import (
	"errors"
	"fmt"
	"os"
	"path"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	ioWriter "github.com/debricked/cli/pkg/io/writer"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	config     conf.IConfig
}

func NewJob(dir string, files []string, cmdFactory ICmdFactory, writer ioWriter.IFileWriter, config conf.IConfig) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(dir, files),
		cmdFactory: cmdFactory,
		config:     config,
	}
}

func (j *Job) Run() {
	fmt.Println("ENTERED RUN")
	workingDirectory := j.GetDir()
	fmt.Println("Files:", j.GetFiles())
	targetClasses := j.GetFiles()[0]
	dependencyDir := ".debrickedTmpFolder"
	targetDir := path.Join(workingDirectory, dependencyDir)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		cmd, _ := j.cmdFactory.MakeBuildMvnCopyDependenciesCmd(workingDirectory, targetDir)
		fmt.Println("building and getting jars", cmd.Args)
		cmd.Output()

	}

	cmd, err := j.cmdFactory.MakeCallGraphGenerationCmd(workingDirectory, targetClasses, targetDir)
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
