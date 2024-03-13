package golang

import (
	"os"
	"syscall"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/io"
	ioFs "github.com/debricked/cli/internal/io"
)

const (
	outputName = "debricked-call-graph-golang"
)

type Job struct {
	job.BaseJob
	mainFile string
	config   conf.IConfig
	archive  io.IArchive
	ctx      cgexec.IContext
	fs       ioFs.IFileSystem
}

func NewJob(dir string, mainFile string, writer ioFs.IFileWriter, archive io.IArchive, config conf.IConfig, ctx cgexec.IContext, fs ioFs.IFileSystem) *Job {
	return &Job{
		BaseJob:  job.NewBaseJob(dir, []string{mainFile}),
		mainFile: mainFile,
		config:   config,
		archive:  archive,
		ctx:      ctx,
		fs:       fs,
	}
}

func (j *Job) Run() {
	workingDirectory := j.GetDir()
	callgraph := NewCallgraphBuilder(
		workingDirectory,
		j.mainFile,
		outputName,
		j.fs,
		j.ctx,
	)
	j.SendStatus("generating call graph")
	j.runCallGraph(&callgraph)
	if j.Errors().HasError() {

		return
	}

}

func (j *Job) runCallGraph(callgraph ICallgraphBuilder) {
	outputFullPath, err := callgraph.RunCallGraph()

	if err != nil {
		j.Errors().Critical(err)

		return
	}
	outputFullPathZip := outputFullPath + ".zip"

	j.SendStatus("zipping callgraph")
	err = j.archive.ZipFile(outputFullPath, outputFullPathZip, outputName)
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
