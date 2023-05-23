package java

import (
	"embed"
	"io"
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/cgexec"
)

//go:embed embeded/SootWrapper.jar
var jarCallGraph embed.FS

type ICallgraph interface {
	RunCallGraphWithSetup() error
	RunCallGraph(callgraphJarPath string) error
}

type Callgraph struct {
	cmdFactory       ICmdFactory
	workingDirectory string
	targetClasses    string
	targetDir        string
	outputName       string
	ctx              cgexec.IContext
}

func (cg *Callgraph) RunCallGraphWithSetup() error {
	jarFile, err := jarCallGraph.Open("embeded/SootWrapper.jar")
	if err != nil {
		return err
	}
	defer jarFile.Close()

	os.TempDir()
	tempDir, err := os.MkdirTemp("", "jar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	tempJarFile := filepath.Join(tempDir, "SootWrapper.jar")

	jarBytes, err := io.ReadAll(jarFile)
	if err != nil {

		return err
	}

	err = os.WriteFile(tempJarFile, jarBytes, 0600)
	if err != nil {

		return err
	}

	err = cg.RunCallGraph(tempJarFile)

	return err
}

func (cg *Callgraph) RunCallGraph(callgraphJarPath string) error {
	cmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(callgraphJarPath, cg.workingDirectory, cg.targetClasses, cg.targetDir, cg.outputName, cg.ctx)
	if err != nil {

		return err
	}

	err = cgexec.RunCommand(cmd, cg.ctx)

	return err
}
