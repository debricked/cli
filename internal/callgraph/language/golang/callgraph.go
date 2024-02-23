package golang

import (
	"github.com/debricked/cli/internal/callgraph/cgexec"
	ioFs "github.com/debricked/cli/internal/io"
)

type ICallgraph interface {
	RunCallGraph() error
}

type Callgraph struct {
	cmdFactory       ICmdFactory
	filesystem       ioFs.IFileSystem
	workingDirectory string
	mainFile         string
	outputName       string
	ctx              cgexec.IContext
}

func NewCallgraph(
	cmdFactory ICmdFactory,
	workingDirectory string,
	mainFile string,
	outputName string,
	filesystem ioFs.IFileSystem,
	ctx cgexec.IContext,
) Callgraph {
	return Callgraph{
		cmdFactory:       cmdFactory,
		workingDirectory: workingDirectory,
		mainFile:         mainFile,
		outputName:       outputName,
		filesystem:       filesystem,
		ctx:              ctx,
	}
}

func (cg *Callgraph) RunCallGraph() error {
	osCmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(cg.mainFile, cg.workingDirectory, cg.ctx)
	if err != nil {
		return err
	}

	cmd := cgexec.NewCommand(osCmd)
	err = cgexec.RunCommand(*cmd, cg.ctx)
	if err != nil {
		return err
	}

	output := cmd.GetStdOut()

	err = cg.filesystem.FsWriteFile(cg.outputName, output.Bytes(), 0600)
	if err != nil {

		return err
	}
	_ = cmd.Wait()

	return err
}
