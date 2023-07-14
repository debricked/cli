package java

import (
	"embed"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	ioFs "github.com/debricked/cli/internal/io"
)

//go:embed embeded/SootWrapper.jar
var jarCallGraph embed.FS

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
	jarFile, err := cg.filesystem.FsOpenEmbed(jarCallGraph, "embeded/SootWrapper.jar")
	if err != nil {
		return err
	}
	defer cg.filesystem.FsCloseFile(jarFile)

	tempDir, err := cg.filesystem.MkdirTemp("jar")
	if err != nil {
		return err
	}
	defer cg.filesystem.RemoveAll(tempDir)
	tempJarFile := filepath.Join(tempDir, "SootWrapper.jar")

	jarBytes, err := cg.filesystem.FsReadAll(jarFile)
	if err != nil {

		return err
	}

	err = cg.filesystem.FsWriteFile(tempJarFile, jarBytes, 0600)
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
