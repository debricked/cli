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
	ioFs "github.com/debricked/cli/internal/io"
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
	fs         ioFs.IFileSystem
}

func NewJob(dir string, files []string, cmdFactory ICmdFactory, writer ioFs.IFileWriter, archive io.IArchive, config conf.IConfig, ctx cgexec.IContext, fs ioFs.IFileSystem) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(dir, files),
		cmdFactory: cmdFactory,
		config:     config,
		archive:    archive,
		ctx:        ctx,
		fs:         fs,
	}
}

func (j *Job) Run() {
	workingDirectory := j.GetDir()
	pmConfig := j.config.PackageManager()
	targetDir := path.Join(workingDirectory, dependencyDir)
	targetClasses := []string{workingDirectory}
	if len(j.GetFiles()) > 0 {
		targetClasses = j.GetFiles()
	}

	// If folder doesn't exist, copy dependencies
	if _, err := j.fs.Stat(targetDir); j.fs.IsNotExist(err) {
		var osCmd *exec.Cmd
		if pmConfig == maven {
			osCmd, err = j.cmdFactory.MakeMvnCopyDependenciesCmd(workingDirectory, targetDir, j.ctx)
			j.SendStatus("copying external dep jars to target folder" + targetDir)
		}
		if err != nil {
			j.Errors().Critical(err)

			return
		}

		j.runCopyDependencies(osCmd)
		if j.Errors().HasError() {
			// If error during copy to .debricked_call_graph, remove the folder

			j.fs.RemoveAll(targetDir)

			return
		}
	}
	callgraph := NewCallgraph(
		j.cmdFactory,
		workingDirectory,
		targetClasses,
		targetDir,
		outputName,
		j.fs,
		j.ctx,
	)
	j.SendStatus("generating call graph")
	j.runCallGraph(&callgraph)
	if j.Errors().HasError() {

		return
	}

	j.runPostProcess()
}

func (j *Job) runCopyDependencies(osCmd *exec.Cmd) {
	cmd := cgexec.NewCommand(osCmd)
	err := cgexec.RunCommand(*cmd, j.ctx)
	if err != nil {
		j.Errors().Critical(err)
	}
}

func (j *Job) runCallGraph(callgraph ICallgraph) {
	err := callgraph.RunCallGraphWithSetup()

	if err != nil {
		j.Errors().Critical(err)

	}
}

func (j *Job) runPostProcess() {
	workingDirectory := j.GetDir()
	outputFullPath := path.Join(workingDirectory, outputName)
	outputFullPathZip := outputFullPath + ".zip"
	j.SendStatus("zipping callgraph")
	err := j.archive.ZipFile(outputFullPath, outputFullPathZip, outputName)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("base64 encoding zipped callgraph")
	err = j.archive.B64(outputFullPathZip, outputFullPath)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("cleanup")
	err = j.archive.Cleanup(outputFullPathZip)
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
