package java

import (
	"github.com/debricked/cli/internal/callgraph/cgexec"
	ioFs "github.com/debricked/cli/internal/io"
)

type ICallgraph interface {
	RunCallGraphWithSetup() error
	RunCallGraph(callgraphJarPath string) error
}

type Callgraph struct {
	cmdFactory       ICmdFactory
	filesystem       ioFs.IFileSystem
	workingDirectory string
	targetClasses    []string
	targetDir        string
	outputName       string
	ctx              cgexec.IContext
}

func NewCallgraph(
	cmdFactory ICmdFactory,
	workingDirectory string,
	targetClasses []string,
	targetDir string,
	outputName string,
	filesystem ioFs.IFileSystem,
	ctx cgexec.IContext,
) Callgraph {
	return Callgraph{
		cmdFactory:       cmdFactory,
		workingDirectory: workingDirectory,
		targetClasses:    targetClasses,
		targetDir:        targetDir,
		outputName:       outputName,
		filesystem:       filesystem,
		ctx:              ctx,
	}
}

func (cg *Callgraph) RunCallGraphWithSetup() error {
	tempDir, err := cg.filesystem.MkdirTemp("jar")
	if err != nil {
		return err
	}
	defer cg.filesystem.RemoveAll(tempDir)

	tempJarFile, err := initializeSootWrapper(cg.filesystem, tempDir)
	if err != nil {

		return err
	}

	err = cg.RunCallGraph(tempJarFile)

	return err
}

func (cg *Callgraph) RunCallGraph(callgraphJarPath string) error {
	osCmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(callgraphJarPath, cg.workingDirectory, cg.targetClasses, cg.targetDir, cg.outputName, cg.ctx)
	if err != nil {

		return err
	}

	cmd := cgexec.NewCommand(osCmd)
	err = cgexec.RunCommand(*cmd, cg.ctx)

	return err
}
