package java

import (
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/io"
	ioWriter "github.com/debricked/cli/internal/io"
)

const (
	maven         = "maven"
	gradle        = "gradle"
	dependencyDir = ".debrickedTmpFolder"
	outputName    = ".debricked-call-graph"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	config     conf.IConfig
	archive    io.IArchive
	ctx        cgexec.IContext
}

func NewJob(dir string, files []string, cmdFactory ICmdFactory, writer ioWriter.IFileWriter, archive io.IArchive, config conf.IConfig, ctx cgexec.IContext) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(dir, files),
		cmdFactory: cmdFactory,
		config:     config,
		archive:    archive,
		ctx:        ctx,
	}
}

func (j *Job) Run() {
	workingDirectory := j.GetDir()
	pmConfig := j.config.Kwargs()["pm"]
	targetDir := path.Join(workingDirectory, dependencyDir)
	targetClasses := workingDirectory
	if len(j.GetFiles()) > 0 {
		targetClasses = j.GetFiles()[0]
	}

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

		j.runCopyDependencies(cmd)

	}
	callgraph := Callgraph{
		cmdFactory:       j.cmdFactory,
		workingDirectory: workingDirectory,
		targetClasses:    targetClasses,
		targetDir:        targetDir,
		outputName:       outputName,
		ctx:              j.ctx,
	}
	j.SendStatus("generating call graph")
	j.runCallGraph(&callgraph)

	j.runPostProcess()
}

func (j *Job) runCopyDependencies(cmd *exec.Cmd) {
	err := cgexec.RunCommand(cmd, j.ctx)
	if err != nil {
		j.Errors().Critical(err)

		return
	}
}

func (j *Job) runCallGraph(callgraph ICallgraph) {
	err := callgraph.RunCallGraphWithSetup()

	if err != nil {
		j.Errors().Critical(err)

		return
	}
}

func (j *Job) runPostProcess() {
	outputNameZip := outputName + ".zip"
	j.SendStatus("zipping callgraph")
	err := j.archive.ZipFile(outputName, outputNameZip)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("base64 encoding zipped callgraph")
	err = j.archive.B64(outputNameZip, outputName)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("cleanup")
	err = j.archive.Cleanup(outputNameZip)
	if err != nil {
		e, ok := err.(*os.PathError)
		if ok && e.Err == syscall.ENOENT {
			return
		} else {
			j.Errors().Critical(err)

			return
		}
	}
}
