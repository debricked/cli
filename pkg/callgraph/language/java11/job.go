package java

import (
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

	// If folder doesn't exist, build
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		cmd, err := j.cmdFactory.MakeBuildMvnCopyDependenciesCmd(workingDirectory, targetDir)
		fmt.Println("building and getting jars", cmd.Args)
		if err != nil {
			j.Errors().Critical(err)

			return
		}
		_, err = cmd.Output()

		if err != nil {
			j.Errors().Critical(err)

			return
		}
	}

	j.SendStatus("generating call graph")
	callgraph := Callgraph{
		cmdFactory:       j.cmdFactory,
		workingDirectory: workingDirectory,
		targetClasses:    targetClasses,
		targetDir:        targetDir,
	}
	err := callgraph.runCallGraphWithSetup()

	if err != nil {
		j.Errors().Critical(err)
	}
}
