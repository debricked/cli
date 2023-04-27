package java

import (
	"os"
	"os/exec"
	"path"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/job"
	ioWriter "github.com/debricked/cli/internal/io/writer"
)

const (
	maven  = "maven"
	gradle = "gradle"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	config     conf.IConfig
	ctx        cgexec.IContext
}

func NewJob(dir string, files []string, cmdFactory ICmdFactory, writer ioWriter.IFileWriter, config conf.IConfig, ctx cgexec.IContext) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(dir, files),
		cmdFactory: cmdFactory,
		config:     config,
		ctx:        ctx,
	}
}

func (j *Job) Run() {
	workingDirectory := j.GetDir()
	targetClasses := j.GetFiles()[0]
	dependencyDir := ".debrickedTmpFolder"
	targetDir := path.Join(workingDirectory, dependencyDir)
	pmConfig := j.config.Kwargs()["pm"]

	// If folder doesn't exist, copy dependencies
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		var cmd *exec.Cmd
		if pmConfig == maven {
			cmd, err = j.cmdFactory.MakeMvnCopyDependenciesCmd(workingDirectory, targetDir, j.ctx)
			j.SendStatus("copying external dep jars to target folder" + targetDir)
		}
		if err != nil {
			j.Errors().Critical(err)

			return
		}

		err = cgexec.RunCommand(cmd, j.ctx)

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
		ctx:              j.ctx,
	}
	err := callgraph.runCallGraphWithSetup()

	if err != nil {
		j.Errors().Critical(err)
	}
}
